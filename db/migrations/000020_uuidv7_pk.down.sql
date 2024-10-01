-- DROP VIEWS
DROP VIEW IF EXISTS v_sos_posts;
DROP VIEW IF EXISTS v_conditions;
DROP VIEW IF EXISTS v_pets_for_sos_posts;
DROP VIEW IF EXISTS v_media_for_sos_posts;

-- DROP INDEXES
DROP INDEX resource_media_resource_id;
DROP INDEX pets_owner_id_idx;
DROP INDEX sos_posts_conditions_sos_post_id;
DROP INDEX sos_posts_pets_sos_post_id;
DROP INDEX sos_posts_dates_sos_post_id;
DROP INDEX sos_posts_author_id_deleted_at;

-- users
-- Add legacy_id column
ALTER TABLE users
    ADD COLUMN legacy_id SERIAL;
ALTER TABLE users
    ADD COLUMN profile_image_legacy_id INTEGER;
-- Rename id to uuid and legacy_id to id
ALTER TABLE users
    DROP CONSTRAINT users_pkey;
ALTER TABLE users
    RENAME COLUMN id TO uuid;
ALTER TABLE users
    RENAME COLUMN legacy_id TO id;
ALTER TABLE users
    ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE users
    ADD PRIMARY KEY (id);
-- profile_image_id -> profile_image_uuid
ALTER TABLE users
    RENAME COLUMN profile_image_id TO profile_image_uuid;
ALTER TABLE users
    RENAME COLUMN profile_image_legacy_id TO profile_image_id;

-- media
-- Add legacy_id column
ALTER TABLE media
    ADD COLUMN legacy_id SERIAL;
-- Rename id to uuid and legacy_id to id
ALTER TABLE media
    DROP CONSTRAINT media_pkey;
ALTER TABLE media
    RENAME COLUMN id TO uuid;
ALTER TABLE media
    RENAME COLUMN legacy_id TO id;
ALTER TABLE media
    ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE media
    ADD PRIMARY KEY (id);

-- breeds
-- Add legacy_id column
ALTER TABLE breeds
    ADD COLUMN legacy_id SERIAL;
-- Rename id to uuid and legacy_id to id
ALTER TABLE breeds
    DROP CONSTRAINT breeds_pkey;
ALTER TABLE breeds
    RENAME COLUMN id TO uuid;
ALTER TABLE breeds
    RENAME COLUMN legacy_id TO id;
ALTER TABLE breeds
    ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE breeds
    ADD PRIMARY KEY (id);

-- pets
-- Add legacy_id column
ALTER TABLE pets
    ADD COLUMN legacy_id SERIAL;
ALTER TABLE pets
    ADD COLUMN owner_legacy_id INTEGER;
ALTER TABLE pets
    ADD COLUMN profile_image_legacy_id INTEGER;
-- Rename id to uuid and legacy_id to id
ALTER TABLE pets
    DROP CONSTRAINT pets_pkey;
ALTER TABLE pets
    RENAME COLUMN id TO uuid;
ALTER TABLE pets
    RENAME COLUMN legacy_id TO id;
ALTER TABLE pets
    ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE pets
    ADD PRIMARY KEY (id);
-- owner_id -> owner_uuid
ALTER TABLE pets
    RENAME COLUMN owner_id TO owner_uuid;
ALTER TABLE pets
    RENAME COLUMN owner_legacy_id TO owner_id;
ALTER TABLE pets
    ALTER COLUMN owner_id DROP NOT NULL;
-- profile_image_id -> profile_image_uuid
ALTER TABLE pets
    RENAME COLUMN profile_image_id TO profile_image_uuid;
ALTER TABLE pets
    RENAME COLUMN profile_image_legacy_id TO profile_image_id;

-- base_posts
CREATE TABLE IF NOT EXISTS base_posts
(
    id         UUID PRIMARY KEY,
    title      VARCHAR(200),
    content    TEXT,
    author_id  UUID,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sos_posts_temp
(
    id           UUID PRIMARY KEY,
    reward       VARCHAR(20),
    care_type    VARCHAR(20),
    carer_gender VARCHAR(10),
    reward_type  VARCHAR(30),
    thumbnail_id UUID
) INHERITS (base_posts);

INSERT INTO sos_posts_temp (id, title, content, author_id, created_at, updated_at, deleted_at)
SELECT id, title, content, author_id, created_at, updated_at, deleted_at
FROM sos_posts;

DROP TABLE sos_posts;

ALTER TABLE sos_posts_temp
    RENAME TO sos_posts;

-- sos_dates
-- Add legacy_id column
ALTER TABLE sos_dates
    ADD COLUMN legacy_id SERIAL;
-- Rename id to uuid and legacy_id to id
ALTER TABLE sos_dates
    DROP CONSTRAINT sos_dates_pkey;
ALTER TABLE sos_dates
    RENAME COLUMN id TO uuid;
ALTER TABLE sos_dates
    RENAME COLUMN legacy_id TO id;
ALTER TABLE sos_dates
    ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE sos_dates
    ADD PRIMARY KEY (id);

-- sos_posts_dates
-- Add legacy_id column
ALTER TABLE sos_posts_dates
    ADD COLUMN legacy_id SERIAL;
ALTER TABLE sos_posts_dates
    ADD COLUMN sos_post_legacy_id INTEGER;
ALTER TABLE sos_posts_dates
    ADD COLUMN sos_dates_legacy_id INTEGER;
-- Rename id to uuid and legacy_id to id
ALTER TABLE sos_posts_dates
    DROP CONSTRAINT sos_posts_dates_pkey;
ALTER TABLE sos_posts_dates
    RENAME COLUMN id TO uuid;
ALTER TABLE sos_posts_dates
    RENAME COLUMN legacy_id TO id;
ALTER TABLE sos_posts_dates
    ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE sos_posts_dates
    ADD PRIMARY KEY (id);
-- sos_post_id -> sos_post_uuid
ALTER TABLE sos_posts_dates
    RENAME COLUMN sos_post_id TO sos_post_uuid;
ALTER TABLE sos_posts_dates
    RENAME COLUMN sos_post_legacy_id TO sos_post_id;
ALTER TABLE sos_posts_dates
    ALTER COLUMN sos_post_uuid DROP NOT NULL;
-- sos_dates_id -> sos_dates_uuid
ALTER TABLE sos_posts_dates
    RENAME COLUMN sos_dates_id TO sos_dates_uuid;
ALTER TABLE sos_posts_dates
    RENAME COLUMN sos_dates_legacy_id TO sos_dates_id;

-- sos_conditions
-- Add legacy_id column
ALTER TABLE sos_conditions
    ADD COLUMN legacy_id SERIAL;
-- Rename id to uuid and legacy_id to id
ALTER TABLE sos_conditions
    DROP CONSTRAINT sos_conditions_pkey;
ALTER TABLE sos_conditions
    RENAME COLUMN id TO uuid;
ALTER TABLE sos_conditions
    RENAME COLUMN legacy_id TO id;
ALTER TABLE sos_conditions
    ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE sos_conditions
    ADD PRIMARY KEY (id);

-- sos_posts_conditions
-- Add legacy_id column
ALTER TABLE sos_posts_conditions
    ADD COLUMN legacy_id SERIAL;
ALTER TABLE sos_posts_conditions
    ADD COLUMN sos_post_legacy_id INTEGER;
ALTER TABLE sos_posts_conditions
    ADD COLUMN sos_condition_legacy_id INTEGER;
-- Rename id to uuid and legacy_id to id
ALTER TABLE sos_posts_conditions
    DROP CONSTRAINT sos_posts_conditions_pkey;
ALTER TABLE sos_posts_conditions
    RENAME COLUMN id TO uuid;
ALTER TABLE sos_posts_conditions
    RENAME COLUMN legacy_id TO id;
ALTER TABLE sos_posts_conditions
    ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE sos_posts_conditions
    ADD PRIMARY KEY (id);
-- sos_post_id -> sos_post_uuid
ALTER TABLE sos_posts_conditions
    RENAME COLUMN sos_post_id TO sos_post_uuid;
ALTER TABLE sos_posts_conditions
    RENAME COLUMN sos_post_legacy_id TO sos_post_id;
ALTER TABLE sos_posts_conditions
    ALTER COLUMN sos_post_uuid DROP NOT NULL;
-- sos_condition_id -> sos_condition_uuid
ALTER TABLE sos_posts_conditions
    RENAME COLUMN sos_condition_id TO sos_condition_uuid;
ALTER TABLE sos_posts_conditions
    RENAME COLUMN sos_condition_legacy_id TO sos_condition_id;

-- sos_posts_pets
-- Add legacy_id column
ALTER TABLE sos_posts_pets
    ADD COLUMN legacy_id SERIAL;
ALTER TABLE sos_posts_pets
    ADD COLUMN sos_post_legacy_id INTEGER;
ALTER TABLE sos_posts_pets
    ADD COLUMN pet_legacy_id INTEGER;
-- Rename id to uuid and legacy_id to id
ALTER TABLE sos_posts_pets
    DROP CONSTRAINT sos_posts_pets_pkey;
ALTER TABLE sos_posts_pets
    RENAME COLUMN id TO uuid;
ALTER TABLE sos_posts_pets
    RENAME COLUMN legacy_id TO id;
ALTER TABLE sos_posts_pets
    ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE sos_posts_pets
    ADD PRIMARY KEY (id);
-- sos_post_id -> sos_post_uuid
ALTER TABLE sos_posts_pets
    RENAME COLUMN sos_post_id TO sos_post_uuid;
ALTER TABLE sos_posts_pets
    RENAME COLUMN sos_post_legacy_id TO sos_post_id;
ALTER TABLE sos_posts_pets
    ALTER COLUMN sos_post_uuid DROP NOT NULL;
-- pet_id -> pet_uuid
ALTER TABLE sos_posts_pets
    RENAME COLUMN pet_id TO pet_uuid;
ALTER TABLE sos_posts_pets
    RENAME COLUMN pet_legacy_id TO pet_id;

-- resource_media
-- Add legacy_id column
ALTER TABLE resource_media
    ADD COLUMN legacy_id SERIAL;
ALTER TABLE resource_media
    ADD COLUMN media_legacy_id INTEGER;
ALTER TABLE resource_media
    ADD COLUMN resource_legacy_id INTEGER;
-- Rename id to uuid and legacy_id to id
ALTER TABLE resource_media
    DROP CONSTRAINT resource_media_pkey;
ALTER TABLE resource_media
    RENAME COLUMN id TO uuid;
ALTER TABLE resource_media
    RENAME COLUMN legacy_id TO id;
ALTER TABLE resource_media
    ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE resource_media
    ADD PRIMARY KEY (id);
-- media_id -> media_uuid
ALTER TABLE resource_media
    RENAME COLUMN media_id TO media_uuid;
ALTER TABLE resource_media
    RENAME COLUMN media_legacy_id TO media_id;
ALTER TABLE resource_media
    ALTER COLUMN media_uuid DROP NOT NULL;
-- resource_id -> resource_uuid
ALTER TABLE resource_media
    RENAME COLUMN resource_id TO resource_uuid;
ALTER TABLE resource_media
    RENAME COLUMN resource_legacy_id TO resource_id;
ALTER TABLE resource_media
    ALTER COLUMN resource_uuid DROP NOT NULL;

-- ADD INDEXES
CREATE INDEX IF NOT EXISTS resource_media_resource_id ON resource_media (resource_id);
CREATE INDEX IF NOT EXISTS pets_owner_id_idx ON pets (owner_id);
CREATE INDEX IF NOT EXISTS sos_posts_conditions_sos_post_id ON sos_posts_conditions (sos_post_id);
CREATE INDEX IF NOT EXISTS sos_posts_pets_sos_post_id ON sos_posts_pets (sos_post_id);
CREATE INDEX IF NOT EXISTS sos_posts_dates_sos_post_id ON sos_posts_dates (sos_post_id);
CREATE INDEX IF NOT EXISTS sos_posts_author_id_deleted_at ON sos_posts (author_id);
