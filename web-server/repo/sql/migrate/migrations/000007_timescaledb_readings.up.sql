-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- readings table uses TEXT PK which blocks hypertable conversion.
-- Hypertables require the partitioning column in the PK.
-- Step 1: Drop the old PK and add timestamp to a composite PK.
ALTER TABLE readings DROP CONSTRAINT readings_pkey;
ALTER TABLE readings ADD PRIMARY KEY (id, timestamp);

-- Step 2: Convert to hypertable partitioned on timestamp.
-- migrate_data => true moves existing rows into chunks.
SELECT create_hypertable('readings', 'timestamp', migrate_data => true);

-- Step 3: Add continuous aggregate for campaign quality metrics.
CREATE MATERIALIZED VIEW campaign_quality_hourly
WITH (timescaledb.continuous) AS
SELECT
    campaign_id,
    time_bucket('1 hour', timestamp) AS bucket,
    count(*) FILTER (WHERE status = 'accepted') AS accepted_count,
    count(*) FILTER (WHERE status = 'quarantined') AS quarantined_count,
    avg(value) FILTER (WHERE status = 'accepted') AS avg_value,
    stddev(value) FILTER (WHERE status = 'accepted') AS stddev_value,
    min(value) FILTER (WHERE status = 'accepted') AS min_value,
    max(value) FILTER (WHERE status = 'accepted') AS max_value
FROM readings
GROUP BY campaign_id, bucket
WITH NO DATA;

-- Refresh policy: keep the aggregate up to date automatically.
SELECT add_continuous_aggregate_policy('campaign_quality_hourly',
    start_offset => INTERVAL '1 day',
    end_offset   => INTERVAL '1 minute',
    schedule_interval => INTERVAL '5 minutes'
);
