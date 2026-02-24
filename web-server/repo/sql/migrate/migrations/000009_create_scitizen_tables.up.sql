-- Scitizen profiles: onboarding state for citizen scientists.
-- No PII stored — only internal references (idp_id lives in app_users).
CREATE TABLE scitizen_profiles (
    user_id         TEXT PRIMARY KEY REFERENCES app_users(id) ON DELETE CASCADE,
    tos_accepted    BOOLEAN NOT NULL DEFAULT false,
    tos_version     TEXT,
    tos_accepted_at TIMESTAMPTZ,
    device_registered   BOOLEAN NOT NULL DEFAULT false,
    campaign_enrolled   BOOLEAN NOT NULL DEFAULT false,
    first_reading       BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Campaign enrollments: links devices to campaigns with consent.
CREATE TABLE campaign_enrollments (
    id              TEXT PRIMARY KEY,
    device_id       TEXT NOT NULL REFERENCES devices(id) ON DELETE CASCADE,
    campaign_id     TEXT NOT NULL REFERENCES campaigns(id) ON DELETE CASCADE,
    scitizen_id     TEXT NOT NULL REFERENCES app_users(id) ON DELETE CASCADE,
    status          TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'withdrawn')),
    enrolled_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    withdrawn_at    TIMESTAMPTZ,
    UNIQUE (device_id, campaign_id)
);

CREATE INDEX idx_enrollments_scitizen ON campaign_enrollments (scitizen_id);
CREATE INDEX idx_enrollments_campaign ON campaign_enrollments (campaign_id);
CREATE INDEX idx_enrollments_device ON campaign_enrollments (device_id);

-- Consent records: immutable audit trail of consent given at enrollment.
CREATE TABLE consent_records (
    id              TEXT PRIMARY KEY,
    enrollment_id   TEXT NOT NULL REFERENCES campaign_enrollments(id) ON DELETE CASCADE,
    version         TEXT NOT NULL,
    scope           TEXT NOT NULL,
    granted_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_consent_enrollment ON consent_records (enrollment_id);

-- Notifications: in-app notification store.
CREATE TABLE notifications (
    id              TEXT PRIMARY KEY,
    user_id         TEXT NOT NULL REFERENCES app_users(id) ON DELETE CASCADE,
    type            TEXT NOT NULL,
    message         TEXT NOT NULL,
    read            BOOLEAN NOT NULL DEFAULT false,
    resource_link   TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_notifications_user ON notifications (user_id, created_at DESC);
CREATE INDEX idx_notifications_unread ON notifications (user_id) WHERE read = false;

-- Notification preferences per user per type.
CREATE TABLE notification_preferences (
    user_id         TEXT NOT NULL REFERENCES app_users(id) ON DELETE CASCADE,
    type            TEXT NOT NULL,
    in_app          BOOLEAN NOT NULL DEFAULT true,
    email           BOOLEAN NOT NULL DEFAULT true,
    PRIMARY KEY (user_id, type)
);
