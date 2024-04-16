ALTER TABLE sos_dates
    ALTER created_at TYPE timestamptz USING created_at AT TIME ZONE 'UTC',
    ALTER updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC',
    ALTER deleted_at TYPE timestamptz USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE sos_posts_dates
    ALTER created_at TYPE timestamptz USING created_at AT TIME ZONE 'UTC',
    ALTER updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC',
    ALTER deleted_at TYPE timestamptz USING deleted_at AT TIME ZONE 'UTC';