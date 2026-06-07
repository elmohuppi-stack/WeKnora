-- Revert source column in knowledges table back to VARCHAR(128)
-- WARNING: This will fail if any rows have source values longer than 128 characters
ALTER TABLE knowledges ALTER COLUMN source TYPE VARCHAR(128);
