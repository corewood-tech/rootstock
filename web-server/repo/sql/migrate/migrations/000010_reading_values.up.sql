-- reading_values: per-parameter measurement values for multi-value readings.
CREATE TABLE reading_values (
    id                TEXT PRIMARY KEY,
    reading_id        TEXT NOT NULL,
    parameter_name    TEXT NOT NULL,
    value             DOUBLE PRECISION NOT NULL,
    status            TEXT NOT NULL DEFAULT 'accepted'
                      CHECK (status IN ('accepted', 'quarantined')),
    quarantine_reason TEXT,
    UNIQUE (reading_id, parameter_name)
);
CREATE INDEX idx_reading_values_reading ON reading_values (reading_id);
CREATE INDEX idx_reading_values_param_status ON reading_values (parameter_name, status);

-- Migrate existing single-value rows into reading_values.
INSERT INTO reading_values (id, reading_id, parameter_name, value, status, quarantine_reason)
SELECT
    gen_random_uuid()::text,
    r.id,
    COALESCE(cp.name, 'value'),
    r.value,
    r.status,
    r.quarantine_reason
FROM readings r
LEFT JOIN LATERAL (
    SELECT name FROM campaign_parameters WHERE campaign_id = r.campaign_id ORDER BY id LIMIT 1
) cp ON true
WHERE r.value IS NOT NULL;

-- readings.value becomes nullable (kept for backward compat during transition).
ALTER TABLE readings ALTER COLUMN value DROP NOT NULL;

-- Replace TimescaleDB continuous aggregate with a regular materialized view
-- that JOINs reading_values (continuous aggregates can't JOIN).
DROP MATERIALIZED VIEW IF EXISTS campaign_quality_hourly;

CREATE MATERIALIZED VIEW campaign_quality_hourly AS
SELECT
    r.campaign_id,
    date_trunc('hour', r.timestamp) AS bucket,
    rv.parameter_name,
    count(*) FILTER (WHERE rv.status = 'accepted') AS accepted_count,
    count(*) FILTER (WHERE rv.status = 'quarantined') AS quarantined_count,
    avg(rv.value) FILTER (WHERE rv.status = 'accepted') AS avg_value,
    stddev(rv.value) FILTER (WHERE rv.status = 'accepted') AS stddev_value,
    min(rv.value) FILTER (WHERE rv.status = 'accepted') AS min_value,
    max(rv.value) FILTER (WHERE rv.status = 'accepted') AS max_value
FROM readings r
JOIN reading_values rv ON rv.reading_id = r.id
GROUP BY r.campaign_id, bucket, rv.parameter_name
WITH NO DATA;

CREATE UNIQUE INDEX idx_cqh_unique ON campaign_quality_hourly (campaign_id, bucket, parameter_name);

-- Populate the materialized view with existing data.
REFRESH MATERIALIZED VIEW campaign_quality_hourly;
