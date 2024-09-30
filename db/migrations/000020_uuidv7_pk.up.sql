-- DROP VIEWS
DROP VIEW IF EXISTS v_sos_posts;
DROP VIEW IF EXISTS v_conditions;
DROP VIEW IF EXISTS v_pets_for_sos_posts;
DROP VIEW IF EXISTS v_media_for_sos_posts;

-- DROP INDEXES
ALTER TABLE resource_media
    DROP CONSTRAINT IF EXISTS resource_media_media_id_fkey;
ALTER TABLE resource_media
    DROP CONSTRAINT IF EXISTS resource_media_resource_id_fkey;
ALTER TABLE sos_posts_conditions
    DROP CONSTRAINT IF EXISTS sos_posts_conditions_sos_post_id_fkey;
ALTER TABLE sos_posts_conditions
    DROP CONSTRAINT IF EXISTS sos_posts_conditions_sos_condition_id_fkey;
ALTER TABLE sos_posts_dates
    DROP CONSTRAINT IF EXISTS sos_posts_dates_sos_dates_id_fkey;
ALTER TABLE sos_posts_dates
    DROP CONSTRAINT IF EXISTS sos_posts_dates_sos_post_id_fkey;
ALTER TABLE sos_posts_pets
    DROP CONSTRAINT IF EXISTS sos_posts_pets_pet_id_fkey;
ALTER TABLE sos_posts_pets
    DROP CONSTRAINT IF EXISTS sos_posts_pets_sos_post_id_fkey;
DROP INDEX resource_media_resource_id;
DROP INDEX pets_owner_id_idx;
DROP INDEX sos_posts_conditions_sos_post_id;
DROP INDEX sos_posts_pets_sos_post_id;
DROP INDEX sos_posts_dates_sos_post_id;
DROP INDEX sos_posts_author_id_deleted_at;

-- users
-- Rename id to legacy_id and uuid to id
ALTER TABLE users
    DROP CONSTRAINT users_pkey;
ALTER TABLE users
    RENAME COLUMN id TO legacy_id;
ALTER TABLE users
    RENAME COLUMN uuid TO id;
ALTER TABLE users
    ALTER COLUMN id SET NOT NULL;
ALTER TABLE users
    ADD PRIMARY KEY (id);
-- profile_image_uuid -> profile_image_id
ALTER TABLE users
    RENAME COLUMN profile_image_id TO profile_image_legacy_id;
ALTER TABLE users
    RENAME COLUMN profile_image_uuid TO profile_image_id;
-- DROP legacy columns
ALTER TABLE users
    DROP COLUMN IF EXISTS legacy_id;
ALTER TABLE users
    DROP COLUMN IF EXISTS profile_image_legacy_id;

-- media
-- Rename id to legacy_id and uuid to id
ALTER TABLE media
    DROP CONSTRAINT media_pkey;
ALTER TABLE media
    RENAME COLUMN id TO legacy_id;
ALTER TABLE media
    RENAME COLUMN uuid TO id;
ALTER TABLE media
    ALTER COLUMN id SET NOT NULL;
ALTER TABLE media
    ADD PRIMARY KEY (id);
-- DROP legacy columns
ALTER TABLE media
    DROP COLUMN IF EXISTS legacy_id;

-- breeds
-- Rename id to legacy_id and uuid to id
ALTER TABLE breeds
    DROP CONSTRAINT breeds_pkey;
ALTER TABLE breeds
    RENAME COLUMN id TO legacy_id;
ALTER TABLE breeds
    RENAME COLUMN uuid TO id;
ALTER TABLE breeds
    ALTER COLUMN id SET NOT NULL;
ALTER TABLE breeds
    ADD PRIMARY KEY (id);
-- DROP legacy columns
ALTER TABLE breeds
    DROP COLUMN IF EXISTS legacy_id;

-- pets
-- Rename id to legacy_id and uuid to id
ALTER TABLE pets
    DROP CONSTRAINT pets_pkey;
ALTER TABLE pets
    RENAME COLUMN id TO legacy_id;
ALTER TABLE pets
    RENAME COLUMN uuid TO id;
ALTER TABLE pets
    ALTER COLUMN id SET NOT NULL;
ALTER TABLE pets
    ADD PRIMARY KEY (id);
-- owner_uuid -> owner_id
ALTER TABLE pets
    RENAME COLUMN owner_id TO owner_legacy_id;
ALTER TABLE pets
    RENAME COLUMN owner_uuid TO owner_id;
ALTER TABLE pets
    ALTER COLUMN owner_id SET NOT NULL;
-- profile_image_uuid -> profile_image_id
ALTER TABLE pets
    RENAME COLUMN profile_image_id TO profile_image_legacy_id;
ALTER TABLE pets
    RENAME COLUMN profile_image_uuid TO profile_image_id;
-- DROP legacy columns
ALTER TABLE pets
    DROP COLUMN IF EXISTS legacy_id;
ALTER TABLE pets
    DROP COLUMN IF EXISTS owner_legacy_id;
ALTER TABLE pets
    DROP COLUMN IF EXISTS profile_image_legacy_id;

-- base_posts
-- Rename id to legacy_id and uuid to id
ALTER TABLE base_posts
    DROP CONSTRAINT base_posts_pkey;
ALTER TABLE base_posts
    RENAME COLUMN id TO legacy_id;
ALTER TABLE base_posts
    RENAME COLUMN uuid TO id;
ALTER TABLE base_posts
    ALTER COLUMN id SET NOT NULL;
ALTER TABLE base_posts
    ADD PRIMARY KEY (id);
-- author_uuid -> author_id
ALTER TABLE base_posts
    RENAME COLUMN author_id TO author_legacy_id;
ALTER TABLE base_posts
    RENAME COLUMN author_uuid TO author_id;
ALTER TABLE base_posts
    ALTER COLUMN author_id SET NOT NULL;
-- DROP legacy columns
ALTER TABLE base_posts
    DROP COLUMN IF EXISTS legacy_id;
ALTER TABLE base_posts
    DROP COLUMN IF EXISTS author_legacy_id;

-- sos_posts
-- thumbnail_uuid -> thumbnail_id
ALTER TABLE sos_posts
    RENAME COLUMN thumbnail_id TO thumbnail_legacy_id;
ALTER TABLE sos_posts
    RENAME COLUMN thumbnail_uuid TO thumbnail_id;
-- DROP legacy columns
ALTER TABLE sos_posts
    DROP COLUMN IF EXISTS thumbnail_legacy_id;

-- sos_dates
-- Rename id to legacy_id and uuid to id
ALTER TABLE sos_dates
    DROP CONSTRAINT sos_dates_pkey;
ALTER TABLE sos_dates
    RENAME COLUMN id TO legacy_id;
ALTER TABLE sos_dates
    RENAME COLUMN uuid TO id;
ALTER TABLE sos_dates
    ALTER COLUMN id SET NOT NULL;
ALTER TABLE sos_dates
    ADD PRIMARY KEY (id);
-- DROP legacy columns
ALTER TABLE sos_dates
    DROP COLUMN IF EXISTS legacy_id;

-- sos_posts_dates
-- Rename id to legacy_id and uuid to id
ALTER TABLE sos_posts_dates
    DROP CONSTRAINT sos_posts_dates_pkey;
ALTER TABLE sos_posts_dates
    RENAME COLUMN id TO legacy_id;
ALTER TABLE sos_posts_dates
    RENAME COLUMN uuid TO id;
ALTER TABLE sos_posts_dates
    ALTER COLUMN id SET NOT NULL;
ALTER TABLE sos_posts_dates
    ADD PRIMARY KEY (id);
-- sos_post_uuid -> sos_post_id
ALTER TABLE sos_posts_dates
    RENAME COLUMN sos_post_id TO sos_post_legacy_id;
ALTER TABLE sos_posts_dates
    RENAME COLUMN sos_post_uuid TO sos_post_id;
ALTER TABLE sos_posts_dates
    ALTER COLUMN sos_post_id SET NOT NULL;
-- sos_dates_uuid -> sos_dates_id
ALTER TABLE sos_posts_dates
    RENAME COLUMN sos_dates_id TO sos_dates_legacy_id;
ALTER TABLE sos_posts_dates
    RENAME COLUMN sos_dates_uuid TO sos_dates_id;
ALTER TABLE sos_posts_dates
    ALTER COLUMN sos_dates_id SET NOT NULL;
-- DROP legacy columns
ALTER TABLE sos_posts_dates
    DROP COLUMN IF EXISTS legacy_id;
ALTER TABLE sos_posts_dates
    DROP COLUMN IF EXISTS sos_post_legacy_id;
ALTER TABLE sos_posts_dates
    DROP COLUMN IF EXISTS sos_dates_legacy_id;

-- sos_conditions
-- Rename id to legacy_id and uuid to id
ALTER TABLE sos_conditions
    DROP CONSTRAINT sos_conditions_pkey;
ALTER TABLE sos_conditions
    RENAME COLUMN id TO legacy_id;
ALTER TABLE sos_conditions
    RENAME COLUMN uuid TO id;
ALTER TABLE sos_conditions
    ALTER COLUMN id SET NOT NULL;
ALTER TABLE sos_conditions
    ADD PRIMARY KEY (id);
-- DROP legacy columns
ALTER TABLE sos_conditions
    DROP COLUMN IF EXISTS legacy_id;

-- sos_posts_conditions
-- Rename id to legacy_id and uuid to id
ALTER TABLE sos_posts_conditions
    DROP CONSTRAINT sos_posts_conditions_pkey;
ALTER TABLE sos_posts_conditions
    RENAME COLUMN id TO legacy_id;
ALTER TABLE sos_posts_conditions
    RENAME COLUMN uuid TO id;
ALTER TABLE sos_posts_conditions
    ALTER COLUMN id SET NOT NULL;
ALTER TABLE sos_posts_conditions
    ADD PRIMARY KEY (id);
-- sos_post_uuid -> sos_post_id
ALTER TABLE sos_posts_conditions
    RENAME COLUMN sos_post_id TO sos_post_legacy_id;
ALTER TABLE sos_posts_conditions
    RENAME COLUMN sos_post_uuid TO sos_post_id;
ALTER TABLE sos_posts_conditions
    ALTER COLUMN sos_post_id SET NOT NULL;
-- sos_condition_uuid -> sos_condition_id
ALTER TABLE sos_posts_conditions
    RENAME COLUMN sos_condition_id TO sos_condition_legacy_id;
ALTER TABLE sos_posts_conditions
    RENAME COLUMN sos_condition_uuid TO sos_condition_id;
-- DROP legacy columns
ALTER TABLE sos_posts_conditions
    DROP COLUMN IF EXISTS legacy_id;
ALTER TABLE sos_posts_conditions
    DROP COLUMN IF EXISTS sos_post_legacy_id;
ALTER TABLE sos_posts_conditions
    DROP COLUMN IF EXISTS sos_condition_legacy_id;

-- sos_posts_pets
-- Rename id to legacy_id and uuid to id
ALTER TABLE sos_posts_pets
    DROP CONSTRAINT sos_posts_pets_pkey;
ALTER TABLE sos_posts_pets
    RENAME COLUMN id TO legacy_id;
ALTER TABLE sos_posts_pets
    RENAME COLUMN uuid TO id;
ALTER TABLE sos_posts_pets
    ALTER COLUMN id SET NOT NULL;
ALTER TABLE sos_posts_pets
    ADD PRIMARY KEY (id);
-- sos_post_uuid -> sos_post_id
ALTER TABLE sos_posts_pets
    RENAME COLUMN sos_post_id TO sos_post_legacy_id;
ALTER TABLE sos_posts_pets
    RENAME COLUMN sos_post_uuid TO sos_post_id;
ALTER TABLE sos_posts_pets
    ALTER COLUMN sos_post_id SET NOT NULL;
-- pet_uuid -> pet_id
ALTER TABLE sos_posts_pets
    RENAME COLUMN pet_id TO pet_legacy_id;
ALTER TABLE sos_posts_pets
    RENAME COLUMN pet_uuid TO pet_id;
ALTER TABLE sos_posts_pets
    ALTER COLUMN pet_id SET NOT NULL;
-- DROP legacy columns
ALTER TABLE sos_posts_pets
    DROP COLUMN IF EXISTS legacy_id;
ALTER TABLE sos_posts_pets
    DROP COLUMN IF EXISTS sos_post_legacy_id;
ALTER TABLE sos_posts_pets
    DROP COLUMN IF EXISTS pet_legacy_id;

-- resource_media
-- Rename id to legacy_id and uuid to id
ALTER TABLE resource_media
    DROP CONSTRAINT resource_media_pkey;
ALTER TABLE resource_media
    RENAME COLUMN id TO legacy_id;
ALTER TABLE resource_media
    RENAME COLUMN uuid TO id;
ALTER TABLE resource_media
    ALTER COLUMN id SET NOT NULL;
ALTER TABLE resource_media
    ADD PRIMARY KEY (id);
-- media_uuid -> media_id
ALTER TABLE resource_media
    RENAME COLUMN media_id TO media_legacy_id;
ALTER TABLE resource_media
    RENAME COLUMN media_uuid TO media_id;
ALTER TABLE resource_media
    ALTER COLUMN media_id SET NOT NULL;
-- resource_uuid -> resource_id
ALTER TABLE resource_media
    RENAME COLUMN resource_id TO resource_legacy_id;
ALTER TABLE resource_media
    RENAME COLUMN resource_uuid TO resource_id;
ALTER TABLE resource_media
    ALTER COLUMN resource_id SET NOT NULL;
-- DROP legacy columns
ALTER TABLE resource_media
    DROP COLUMN IF EXISTS legacy_id;
ALTER TABLE resource_media
    DROP COLUMN IF EXISTS media_legacy_id;
ALTER TABLE resource_media
    DROP COLUMN IF EXISTS resource_legacy_id;

-- ADD INDEXES
CREATE INDEX IF NOT EXISTS resource_media_resource_id ON resource_media (resource_id);
CREATE INDEX IF NOT EXISTS pets_owner_id_idx ON pets (owner_id);
CREATE INDEX IF NOT EXISTS sos_posts_conditions_sos_post_id ON sos_posts_conditions (sos_post_id);
CREATE INDEX IF NOT EXISTS sos_posts_pets_sos_post_id ON sos_posts_pets (sos_post_id);
CREATE INDEX IF NOT EXISTS sos_posts_dates_sos_post_id ON sos_posts_dates (sos_post_id);
CREATE INDEX IF NOT EXISTS sos_posts_author_id_deleted_at ON sos_posts (author_id);
