CREATE TABLE IF NOT EXISTS base_posts
(
    id         SERIAL PRIMARY KEY,
    title      VARCHAR(200),
    content    TEXT,
    author_id  BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);