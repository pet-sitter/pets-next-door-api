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
-- Add legacy_id column
ALTER TABLE base_posts
    ADD COLUMN legacy_id SERIAL;
ALTER TABLE base_posts
    ADD COLUMN author_legacy_id INTEGER;
-- Rename id to uuid and legacy_id to id
ALTER TABLE base_posts
    DROP CONSTRAINT base_posts_pkey;
ALTER TABLE base_posts
    RENAME COLUMN id TO uuid;
ALTER TABLE base_posts
    RENAME COLUMN legacy_id TO id;
ALTER TABLE base_posts
    ALTER COLUMN uuid SET NOT NULL;
ALTER TABLE base_posts
    ADD PRIMARY KEY (id);
-- author_id -> author_uuid
ALTER TABLE base_posts
    RENAME COLUMN author_id TO author_uuid;
ALTER TABLE base_posts
    RENAME COLUMN author_legacy_id TO author_id;
ALTER TABLE base_posts
    ALTER COLUMN author_uuid DROP NOT NULL;

-- sos_posts
-- Add legacy_id column
ALTER TABLE sos_posts
    ADD COLUMN thumbnail_legacy_id INTEGER;
-- thumbnail_id -> thumbnail_uuid
ALTER TABLE sos_posts
    RENAME COLUMN thumbnail_id TO thumbnail_uuid;
ALTER TABLE sos_posts
    RENAME COLUMN thumbnail_legacy_id TO thumbnail_id;

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

-- Add view
-- 돌봄 급구(SosPosts) 테이블 VIEW 생성
CREATE OR REPLACE VIEW v_sos_posts AS
SELECT sos_posts.id,
       sos_posts.title,
       sos_posts.content,
       sos_posts.reward,
       sos_posts.reward_type,
       sos_posts.care_type,
       sos_posts.carer_gender,
       sos_posts.thumbnail_id,
       sos_posts.author_id,
       sos_posts.created_at,
       sos_posts.updated_at,
       MIN(sos_dates.date_start_at)                                      AS earliest_date_start_at,
       json_agg(sos_dates.*) FILTER (WHERE sos_dates.deleted_at IS NULL) AS dates
FROM sos_posts
         LEFT JOIN sos_posts_dates ON sos_posts.id = sos_posts_dates.sos_post_id
         LEFT JOIN sos_dates ON sos_posts_dates.sos_dates_id = sos_dates.id
WHERE sos_posts.deleted_at IS NULL
  AND sos_dates.deleted_at IS NULL
  AND sos_posts_dates.deleted_at IS NULL
GROUP BY sos_posts.id;

-- 돌봄 급구 Conditions 테이블 VIEW 생성
CREATE OR REPLACE VIEW v_conditions AS
SELECT sos_posts_conditions.sos_post_id,
       json_agg(sos_conditions.*)
       FILTER (WHERE sos_conditions.deleted_at IS NULL) AS conditions_info
FROM sos_posts_conditions
         LEFT JOIN sos_conditions ON sos_posts_conditions.sos_condition_id = sos_conditions.id
WHERE sos_conditions.deleted_at IS NULL
  AND sos_posts_conditions.deleted_at IS NULL
GROUP BY sos_posts_conditions.sos_post_id;

-- 돌봄 급구 관련 Pets 테이블 VIEW 생성
CREATE OR REPLACE VIEW v_pets_for_sos_posts AS
SELECT sos_posts_pets.sos_post_id,
       array_agg(pets.pet_type)                         AS pet_type_list,
       json_agg(
       json_build_object(
               'id', pets.id,
               'owner_id', pets.owner_id,
               'name', pets.name,
               'pet_type', pets.pet_type,
               'sex', pets.sex,
               'neutered', pets.neutered,
               'breed', pets.breed,
               'birth_date', pets.birth_date,
               'weight_in_kg', pets.weight_in_kg,
               'additional_note', pets.additional_note,
               'created_at', pets.created_at,
               'updated_at', pets.updated_at,
               'deleted_at', pets.deleted_at,
               'remarks', pets.remarks,
               'profile_image_id', pets.profile_image_id,
               'profile_image_url', media.url
       )
               ) FILTER (WHERE pets.deleted_at IS NULL) AS pets_info
FROM sos_posts_pets
         INNER JOIN pets ON sos_posts_pets.pet_id = pets.id AND pets.deleted_at IS NULL
         LEFT JOIN media ON pets.profile_image_id = media.id
WHERE sos_posts_pets.deleted_at IS NULL
GROUP BY sos_posts_pets.sos_post_id;



-- 돌봄 급구 관련 Media 테이블 VIEW 생성
CREATE OR REPLACE VIEW v_media_for_sos_posts AS
SELECT resource_media.resource_id                                AS sos_post_id,
       json_agg(media.*) FILTER (WHERE media.deleted_at IS NULL) AS media_info
FROM resource_media
         LEFT JOIN media ON resource_media.media_id = media.id
WHERE media.deleted_at IS NULL
  AND resource_media.deleted_at IS NULL
GROUP BY resource_media.resource_id;
