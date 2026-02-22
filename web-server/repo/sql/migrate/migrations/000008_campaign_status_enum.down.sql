-- Revert to original 3-state constraint.
UPDATE campaigns SET status = 'closed' WHERE status IN ('cancelled', 'active', 'completed');
ALTER TABLE campaigns DROP CONSTRAINT IF EXISTS campaigns_status_check;
ALTER TABLE campaigns ADD CONSTRAINT campaigns_status_check
    CHECK (status IN ('draft', 'published', 'closed'));
