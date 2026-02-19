CREATE TABLE readings (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    device_id           UUID            NOT NULL REFERENCES devices(id),
    campaign_id         UUID            NOT NULL REFERENCES campaigns(id),
    value               DOUBLE PRECISION NOT NULL,
    timestamp           TIMESTAMPTZ     NOT NULL,
    geolocation         JSONB,
    firmware_version    TEXT            NOT NULL,
    cert_serial         TEXT            NOT NULL,
    ingested_at         TIMESTAMPTZ     NOT NULL DEFAULT now(),
    status              TEXT            NOT NULL DEFAULT 'accepted' CHECK (status IN ('accepted', 'quarantined')),
    quarantine_reason   TEXT
);

CREATE INDEX idx_readings_device_campaign ON readings (device_id, campaign_id);
CREATE INDEX idx_readings_campaign_status ON readings (campaign_id, status);
CREATE INDEX idx_readings_ingested_at ON readings (ingested_at);
