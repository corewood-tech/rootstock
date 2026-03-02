import pg from 'pg';
import { randomBytes } from 'node:crypto';

const { Pool } = pg;

let pool: pg.Pool | null = null;

/** Singleton Postgres pool connected to the app database. */
export function getPool(): pg.Pool {
  if (!pool) {
    pool = new Pool({
      host: 'app-postgres',
      port: 5432,
      user: 'rootstock',
      password: 'rootstock',
      database: 'rootstock',
    });
  }
  return pool;
}

/** Look up the internal app_users.id for a given IdP subject. */
export async function getUserID(idpID: string): Promise<string> {
  const result = await getPool().query(
    'SELECT id FROM app_users WHERE idp_id = $1',
    [idpID],
  );
  if (result.rows.length === 0) {
    throw new Error(`No app_user found for idp_id: ${idpID}`);
  }
  return result.rows[0].id;
}

/** Look up the internal user ID by email-pattern match (test users). */
export async function getUserIDByPattern(emailPattern: string): Promise<string> {
  const result = await getPool().query(
    'SELECT id FROM app_users ORDER BY created_at DESC LIMIT 1',
  );
  if (result.rows.length === 0) {
    throw new Error('No app_users found');
  }
  return result.rows[0].id;
}

/** Get the most recent campaign ID, optionally filtered by status. */
export async function getCampaignID(status?: string): Promise<string> {
  const query = status
    ? 'SELECT id FROM campaigns WHERE status = $1 ORDER BY created_at DESC LIMIT 1'
    : 'SELECT id FROM campaigns ORDER BY created_at DESC LIMIT 1';
  const params = status ? [status] : [];
  const result = await getPool().query(query, params);
  if (result.rows.length === 0) {
    throw new Error(`No campaign found${status ? ` with status: ${status}` : ''}`);
  }
  return result.rows[0].id;
}

/** Get the creator (app_users.id) of a campaign. */
export async function getCampaignCreator(campaignID: string): Promise<string> {
  const result = await getPool().query(
    'SELECT created_by FROM campaigns WHERE id = $1',
    [campaignID],
  );
  if (result.rows.length === 0) {
    throw new Error(`No campaign found with id: ${campaignID}`);
  }
  return result.rows[0].created_by;
}

/** Generate a ULID-style random ID (26 chars, Crockford base32). */
function generateID(): string {
  const ENCODING = '0123456789ABCDEFGHJKMNPQRSTVWXYZ';
  const bytes = randomBytes(16);
  let id = '';
  for (const b of bytes) {
    id += ENCODING[b % 32];
  }
  return id.slice(0, 26);
}

/** Insert a device into the devices table. Returns the device ID. */
export async function seedDevice(
  ownerID: string,
  deviceClass: string,
  firmwareVersion: string,
  tier: number,
  sensors: string[],
): Promise<string> {
  const id = generateID();
  await getPool().query(
    `INSERT INTO devices (id, owner_id, status, class, firmware_version, tier, sensors)
     VALUES ($1, $2, 'pending', $3, $4, $5, $6)`,
    [id, ownerID, deviceClass, firmwareVersion, tier, sensors],
  );
  return id;
}

/** Insert an enrollment code for a device. */
export async function seedEnrollmentCode(
  deviceID: string,
  code: string,
  ttlSeconds: number,
): Promise<void> {
  const expiresAt = new Date(Date.now() + ttlSeconds * 1000);
  await getPool().query(
    `INSERT INTO enrollment_codes (code, device_id, expires_at, used)
     VALUES ($1, $2, $3, false)`,
    [code, deviceID, expiresAt],
  );
}

/** Insert a campaign enrollment linking device → campaign → scitizen. */
export async function seedEnrollment(
  deviceID: string,
  campaignID: string,
  scitizenID: string,
): Promise<string> {
  const id = generateID();
  await getPool().query(
    `INSERT INTO campaign_enrollments (id, device_id, campaign_id, scitizen_id, status, enrolled_at)
     VALUES ($1, $2, $3, $4, 'active', now())`,
    [id, deviceID, campaignID, scitizenID],
  );
  return id;
}

/** Promote a user to 'both' role so they can access scitizen endpoints. */
export async function promoteUserToBoth(userID: string): Promise<void> {
  await getPool().query(
    `UPDATE app_users SET user_type = 'both' WHERE id = $1`,
    [userID],
  );
}

/** Seed a scitizen_profiles row (required for onboarding state). */
export async function seedScitizenProfile(userID: string): Promise<void> {
  await getPool().query(
    `INSERT INTO scitizen_profiles (user_id, tos_accepted, device_registered, campaign_enrolled, first_reading)
     VALUES ($1, true, true, true, false)
     ON CONFLICT (user_id) DO NOTHING`,
    [userID],
  );
}

/** Close the pool. Call in afterAll. */
export async function cleanup(): Promise<void> {
  if (pool) {
    await pool.end();
    pool = null;
  }
}
