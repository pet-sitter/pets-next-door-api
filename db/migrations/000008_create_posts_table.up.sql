CREATE TABLE IF NOT EXISTS sos_posts (
 id SERIAL PRIMARY KEY,
 reward VARCHAR(20),
 date_start_at DATE,
 date_end_at DATE,
 time_start_at TIME,
 time_end_at TIME,
 care_type VARCHAR(20),
 carer_gender VARCHAR(10),
 reward_amount VARCHAR(30),
 thumbnail_ID BIGINT
) INHERITS (base_posts);

CREATE INDEX IF NOT EXISTS sos_posts_author_id_deleted_at ON sos_posts(author_id);

CREATE TABLE IF NOT EXISTS resource_media (
  id SERIAL PRIMARY KEY,
  media_id BIGINT REFERENCES media(id),
  resource_id BIGINT REFERENCES sos_posts(id),
  resource_type VARCHAR(20),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS resource_media_resource_id ON resource_media(resource_id);

CREATE TABLE IF NOT EXISTS sos_conditions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sos_posts_conditions (
  id SERIAL PRIMARY KEY,
  sos_post_id BIGINT REFERENCES sos_posts(id),
  sos_condition_id BIGINT REFERENCES sos_conditions(id),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS sos_posts_conditions_sos_post_id ON sos_posts_conditions(sos_post_id);

CREATE TABLE IF NOT EXISTS sos_posts_pets (
  id SERIAL PRIMARY KEY,
  sos_post_id BIGINT REFERENCES sos_posts(id),
  pet_id BIGINT REFERENCES pets(id),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS sos_posts_pets_sos_post_id ON sos_posts_pets(sos_post_id);
