-- Remove continuous aggregate policy and view
DROP MATERIALIZED VIEW IF EXISTS campaign_quality_hourly;

-- Revert readings from hypertable back to regular table.
-- TimescaleDB does not support reverting a hypertable in-place.
-- Recreate the table from the hypertable data.
CREATE TABLE readings_plain AS SELECT * FROM readings;
DROP TABLE readings;
ALTER TABLE readings_plain RENAME TO readings;
ALTER TABLE readings ADD PRIMARY KEY (id);
ALTER TABLE readings ADD CONSTRAINT readings_device_id_fkey FOREIGN KEY (device_id) REFERENCES devices(id);
ALTER TABLE readings ADD CONSTRAINT readings_campaign_id_fkey FOREIGN KEY (campaign_id) REFERENCES campaigns(id);
CREATE INDEX idx_readings_device_campaign ON readings (device_id, campaign_id);
CREATE INDEX idx_readings_campaign_status ON readings (campaign_id, status);
CREATE INDEX idx_readings_ingested_at ON readings (ingested_at);

DROP EXTENSION IF EXISTS timescaledb;
