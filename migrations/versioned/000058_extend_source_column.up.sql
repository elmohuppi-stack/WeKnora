-- Extend source column in knowledges table from VARCHAR(128) to TEXT
-- to support long URLs (e.g. https://example.com/very/long/path/...)
ALTER TABLE knowledges ALTER COLUMN source TYPE TEXT;
