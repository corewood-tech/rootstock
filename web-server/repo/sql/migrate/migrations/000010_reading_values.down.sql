-- Revert materialized view back to TimescaleDB continuous aggregate.
DROP MATERIALIZED VIEW IF EXISTS campaign_quality_hourly;

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

SELECT add_continuous_aggregate_policy('campaign_quality_hourly',
    start_offset => INTERVAL '1 day',
    end_offset   => INTERVAL '1 minute',
    schedule_interval => INTERVAL '5 minutes'
);

-- Restore readings.value NOT NULL from reading_values data.
UPDATE readings r
SET value = (
    SELECT rv.value FROM reading_values rv WHERE rv.reading_id = r.id ORDER BY rv.id LIMIT 1
)
WHERE r.value IS NULL;

ALTER TABLE readings ALTER COLUMN value SET NOT NULL;

-- Drop reading_values table.
DROP TABLE IF EXISTS reading_values;
