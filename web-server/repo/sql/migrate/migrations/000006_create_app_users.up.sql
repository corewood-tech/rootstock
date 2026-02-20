CREATE TABLE app_users (
    id          TEXT PRIMARY KEY,
    idp_id      TEXT NOT NULL UNIQUE,
    user_type   TEXT NOT NULL CHECK (user_type IN ('scitizen', 'researcher', 'both')),
    status      TEXT NOT NULL DEFAULT 'active',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_app_users_idp_id ON app_users (idp_id);
