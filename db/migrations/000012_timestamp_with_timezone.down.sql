ALTER TABLE users
    ALTER created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC',
    ALTER updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC',
    ALTER deleted_at TYPE timestamp USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE media
    ALTER created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC',
    ALTER updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC',
    ALTER deleted_at TYPE timestamp USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE pets
    ALTER created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC',
    ALTER updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC',
    ALTER deleted_at TYPE timestamp USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE breeds
    ALTER created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC',
    ALTER updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC',
    ALTER deleted_at TYPE timestamp USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE base_posts
    ALTER created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC',
    ALTER updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC',
    ALTER deleted_at TYPE timestamp USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE resource_media
    ALTER created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC',
    ALTER updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC',
    ALTER deleted_at TYPE timestamp USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE sos_conditions
    ALTER created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC',
    ALTER updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC',
    ALTER deleted_at TYPE timestamp USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE sos_posts_conditions
    ALTER created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC',
    ALTER updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC',
    ALTER deleted_at TYPE timestamp USING deleted_at AT TIME ZONE 'UTC';

ALTER TABLE sos_posts_pets
    ALTER created_at TYPE timestamp USING created_at AT TIME ZONE 'UTC',
    ALTER updated_at TYPE timestamp USING updated_at AT TIME ZONE 'UTC',
    ALTER deleted_at TYPE timestamp USING deleted_at AT TIME ZONE 'UTC';
