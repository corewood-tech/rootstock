CREATE TABLE devices (
    id                  TEXT PRIMARY KEY,
    owner_id            TEXT        NOT NULL,
    status              TEXT        NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'active', 'suspended', 'revoked')),
    class               TEXT        NOT NULL,
    firmware_version    TEXT        NOT NULL,
    tier                INT         NOT NULL,
    sensors             TEXT[]      NOT NULL DEFAULT '{}',
    cert_serial         TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE enrollment_codes (
    code        TEXT PRIMARY KEY,
    device_id   TEXT        NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    expires_at  TIMESTAMPTZ NOT NULL,
    used        BOOLEAN     NOT NULL DEFAULT false
);

CREATE TABLE device_campaigns (
    device_id   TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    campaign_id TEXT NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (device_id, campaign_id)
);
