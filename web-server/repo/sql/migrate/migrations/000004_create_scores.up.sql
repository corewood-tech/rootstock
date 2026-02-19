CREATE TABLE scores (
    scitizen_id     TEXT PRIMARY KEY,
    volume          INT             NOT NULL DEFAULT 0,
    quality_rate    DOUBLE PRECISION NOT NULL DEFAULT 0,
    consistency     DOUBLE PRECISION NOT NULL DEFAULT 0,
    diversity       INT             NOT NULL DEFAULT 0,
    total           DOUBLE PRECISION NOT NULL DEFAULT 0,
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT now()
);

CREATE TABLE badges (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scitizen_id TEXT        NOT NULL,
    badge_type  TEXT        NOT NULL,
    awarded_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE sweepstakes_entries (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    scitizen_id         TEXT        NOT NULL,
    entries             INT         NOT NULL,
    milestone_trigger   TEXT        NOT NULL,
    granted_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
