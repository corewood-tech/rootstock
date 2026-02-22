-- Align campaign status values with lifecycle spec (graph 0x26-0x2a).
-- Add: active, completed, cancelled. Remove: closed (replaced by cancelled).
ALTER TABLE campaigns DROP CONSTRAINT IF EXISTS campaigns_status_check;
ALTER TABLE campaigns ADD CONSTRAINT campaigns_status_check
    CHECK (status IN ('draft', 'published', 'active', 'completed', 'cancelled'));

-- Migrate any existing 'closed' rows to 'cancelled'.
UPDATE campaigns SET status = 'cancelled' WHERE status = 'closed';
