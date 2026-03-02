import { execSync } from 'node:child_process';
import { mkdtempSync, readFileSync, rmSync } from 'node:fs';
import { join } from 'node:path';
import { tmpdir } from 'node:os';
import { randomBytes } from 'node:crypto';
import mqtt from 'mqtt';
import type { APIRequestContext } from '@playwright/test';
import * as db from './db';

/**
 * MockDevice simulates a physical IoT device for E2E testing.
 * Handles enrollment (CSR generation → cert issuance) and MQTT publishing with mTLS.
 */
export class MockDevice {
  deviceId: string;
  certPEM: string = '';
  serial: string = '';
  firmwareVersion: string;

  private privateKeyPEM: string;
  private mqttClient: mqtt.MqttClient | null = null;
  private tmpDir: string;

  private constructor(
    deviceId: string,
    privateKeyPEM: string,
    firmwareVersion: string,
    tmpDir: string,
  ) {
    this.deviceId = deviceId;
    this.privateKeyPEM = privateKeyPEM;
    this.firmwareVersion = firmwareVersion;
    this.tmpDir = tmpDir;
  }

  /**
   * Seed a device in the DB, generate an enrollment code, and enroll via HTTP.
   * Returns a fully enrolled MockDevice with a signed certificate.
   */
  static async create(
    request: APIRequestContext,
    ownerID: string,
    opts?: {
      deviceClass?: string;
      firmwareVersion?: string;
      tier?: number;
      sensors?: string[];
    },
  ): Promise<MockDevice> {
    const deviceClass = opts?.deviceClass ?? 'air-quality-monitor';
    const firmwareVersion = opts?.firmwareVersion ?? '1.0.0';
    const tier = opts?.tier ?? 1;
    const sensors = opts?.sensors ?? ['PM2.5'];

    // Seed device in DB
    const deviceId = await db.seedDevice(
      ownerID,
      deviceClass,
      firmwareVersion,
      tier,
      sensors,
    );

    // Seed enrollment code (valid for 5 minutes)
    const enrollmentCode = randomBytes(16).toString('hex');
    await db.seedEnrollmentCode(deviceId, enrollmentCode, 300);

    // Generate EC P-256 key pair and CSR via openssl
    const tmpDir = mkdtempSync(join(tmpdir(), 'mock-device-'));
    const keyPath = join(tmpDir, 'device.key');
    const csrPath = join(tmpDir, 'device.csr');

    execSync(
      `openssl ecparam -genkey -name prime256v1 -noout -out ${keyPath}`,
      { stdio: 'pipe' },
    );
    execSync(
      `openssl req -new -key ${keyPath} -out ${csrPath} -subj "/CN=${deviceId}"`,
      { stdio: 'pipe' },
    );

    const privateKeyPEM = readFileSync(keyPath, 'utf-8');
    const csrPEM = readFileSync(csrPath, 'utf-8');

    const device = new MockDevice(deviceId, privateKeyPEM, firmwareVersion, tmpDir);

    // Enroll via HTTP
    await device.enroll(request, enrollmentCode, csrPEM);

    return device;
  }

  /** POST /enroll with the enrollment code and CSR. Stores the issued certificate. */
  async enroll(
    request: APIRequestContext,
    enrollmentCode: string,
    csrPEM: string,
  ): Promise<void> {
    const response = await request.post('http://caddy:9999/enroll', {
      data: {
        enrollment_code: enrollmentCode,
        csr: csrPEM,
      },
    });

    if (!response.ok()) {
      const body = await response.text();
      throw new Error(`Enrollment failed (${response.status()}): ${body}`);
    }

    const result = await response.json();
    this.certPEM = result.cert_pem;
    this.serial = result.serial;
  }

  /** Publish a reading to the MQTT broker using mTLS. */
  async publishReading(
    request: APIRequestContext,
    campaignID: string,
    value: number,
  ): Promise<void> {
    // Fetch CA cert for TLS verification
    const caResponse = await request.get('http://caddy:9999/ca');
    if (!caResponse.ok()) {
      throw new Error(`Failed to fetch CA cert: ${caResponse.status()}`);
    }
    const caCertPEM = await caResponse.text();

    // Connect via MQTT with mTLS — disable auto-reconnect so auth failures surface
    const client = await mqtt.connectAsync(
      'mqtts://web-server:8883',
      {
        clientId: this.deviceId,
        key: this.privateKeyPEM,
        cert: this.certPEM,
        ca: caCertPEM,
        rejectUnauthorized: true,
        protocolVersion: 5,
        reconnectPeriod: 0,
      },
      false, // reject on first failure
    );
    this.mqttClient = client;

    const topic = `rootstock/${this.deviceId}/data/${campaignID}`;
    const payload = JSON.stringify({
      value,
      timestamp: new Date().toISOString(),
      firmware_version: this.firmwareVersion,
      cert_serial: this.serial,
    });

    await client.publishAsync(topic, payload, { qos: 1 });
  }

  /** Disconnect the MQTT client and clean up temp files. */
  async disconnect(): Promise<void> {
    if (this.mqttClient) {
      await this.mqttClient.endAsync();
      this.mqttClient = null;
    }
    try {
      rmSync(this.tmpDir, { recursive: true, force: true });
    } catch {
      // best-effort cleanup
    }
  }
}
