-- 002: Add columns added after initial schema
-- Safe to re-run (uses IF NOT EXISTS pattern via IGNORE)

ALTER TABLE tenants ADD COLUMN max_emails_per_month INT DEFAULT NULL;
ALTER TABLE recipients ADD COLUMN gender VARCHAR(10) DEFAULT '' AFTER department;
ALTER TABLE campaigns ADD COLUMN schedule_start DATETIME DEFAULT NULL AFTER send_by;
ALTER TABLE campaigns ADD COLUMN working_hours_only BOOLEAN NOT NULL DEFAULT FALSE AFTER schedule_start;
ALTER TABLE campaigns ADD COLUMN skip_weekends BOOLEAN NOT NULL DEFAULT FALSE AFTER working_hours_only;
ALTER TABLE campaigns ADD COLUMN selection_mode VARCHAR(20) NOT NULL DEFAULT 'all' AFTER phish_url;
ALTER TABLE campaigns ADD COLUMN sample_percent INT NOT NULL DEFAULT 100 AFTER selection_mode;
ALTER TABLE campaigns ADD COLUMN departments TEXT AFTER sample_percent;
