ALTER TABLE badges ADD CONSTRAINT badges_scitizen_badge_type_unique UNIQUE (scitizen_id, badge_type);
ALTER TABLE sweepstakes_entries ADD CONSTRAINT sweepstakes_scitizen_milestone_unique UNIQUE (scitizen_id, milestone_trigger);
