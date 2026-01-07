ALTER TABLE transactions ADD COLUMN IF NOT EXISTS scope String DEFAULT '{}';
