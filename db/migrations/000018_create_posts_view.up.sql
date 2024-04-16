-- 돌봄 급구(SosPosts) 테이블 VIEW 생성
CREATE OR REPLACE VIEW v_sos_posts AS
SELECT
    sos_posts.id,
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
    MIN(sos_dates.date_start_at) AS earliest_date_start_at,
    json_agg(sos_dates.*) FILTER (WHERE sos_dates.deleted_at IS NULL) AS dates
FROM
    sos_posts
        LEFT JOIN sos_posts_dates ON sos_posts.id = sos_posts_dates.sos_post_id
        LEFT JOIN sos_dates ON sos_posts_dates.sos_dates_id = sos_dates.id
WHERE
    sos_posts.deleted_at IS NULL
  AND sos_dates.deleted_at IS NULL
  AND sos_posts_dates.deleted_at IS NULL
GROUP BY sos_posts.id;

-- 돌봄 급구 Conditions 테이블 VIEW 생성
CREATE OR REPLACE VIEW v_conditions AS
SELECT
    sos_posts_conditions.sos_post_id,
    json_agg(sos_conditions.*) FILTER (WHERE sos_conditions.deleted_at IS NULL) AS conditions_info
FROM
    sos_posts_conditions
        LEFT JOIN sos_conditions ON sos_posts_conditions.sos_condition_id = sos_conditions.id
WHERE
    sos_conditions.deleted_at IS NULL AND
    sos_posts_conditions.deleted_at IS NULL
GROUP BY sos_posts_conditions.sos_post_id;

-- 돌봄 급구 관련 Pets 테이블 VIEW 생성
CREATE OR REPLACE VIEW v_pets_for_sos_posts AS
SELECT
    sos_posts_pets.sos_post_id,
    array_agg(pets.pet_type) AS pet_type_list,
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
FROM
    sos_posts_pets
        INNER JOIN pets ON sos_posts_pets.pet_id = pets.id AND pets.deleted_at IS NULL
        LEFT JOIN media ON pets.profile_image_id = media.id
WHERE
    sos_posts_pets.deleted_at IS NULL
GROUP BY sos_posts_pets.sos_post_id;



-- 돌봄 급구 관련 Media 테이블 VIEW 생성
CREATE OR REPLACE VIEW v_media_for_sos_posts AS
SELECT
    resource_media.resource_id AS sos_post_id,
    json_agg(media.*) FILTER (WHERE media.deleted_at IS NULL) AS media_info
FROM
    resource_media
        LEFT JOIN media ON resource_media.media_id = media.id
WHERE
    media.deleted_at IS NULL AND
    resource_media.deleted_at IS NULL
GROUP BY resource_media.resource_id;