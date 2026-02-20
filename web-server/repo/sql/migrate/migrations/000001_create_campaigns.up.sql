CREATE TABLE campaigns (
    id              TEXT PRIMARY KEY,
    org_id          TEXT        NOT NULL,
    status          TEXT        NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'closed')),
    window_start    TIMESTAMPTZ,
    window_end      TIMESTAMPTZ,
    created_by      TEXT        NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE campaign_parameters (
    id              TEXT PRIMARY KEY,
    campaign_id     TEXT        NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    name            TEXT        NOT NULL,
    unit            TEXT        NOT NULL,
    min_range       DOUBLE PRECISION,
    max_range       DOUBLE PRECISION,
    precision       INT
);

CREATE TABLE campaign_regions (
    id              TEXT PRIMARY KEY,
    campaign_id     TEXT        NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    geometry        JSONB       NOT NULL
);

CREATE TABLE campaign_eligibility (
    id              TEXT PRIMARY KEY,
    campaign_id     TEXT        NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    device_class    TEXT        NOT NULL,
    tier            INT         NOT NULL,
    required_sensors TEXT[]     NOT NULL DEFAULT '{}',
    firmware_min    TEXT
);
