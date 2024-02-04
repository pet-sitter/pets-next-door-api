ALTER TABLE sos_posts
    DROP COLUMN IF EXISTS time_start_at,
    DROP COLUMN IF EXISTS time_end_at;