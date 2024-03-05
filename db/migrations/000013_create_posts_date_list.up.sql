ALTER TABLE sos_posts
    DROP COLUMN IF EXISTS date_start_at,
    DROP COLUMN IF EXISTS date_end_at;

CREATE TABLE IF NOT EXISTS sos_dates
(
    id            SERIAL PRIMARY KEY,
    date_start_at DATE,
    date_end_at   DATE,
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS sos_posts_dates
(
    id            SERIAL PRIMARY KEY,
    sos_post_id   BIGINT REFERENCES sos_posts (id),
    sos_dates_id  BIGINT REFERENCES sos_dates (id),
    created_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS sos_posts_dates_sos_post_id ON sos_posts_dates (sos_post_id);