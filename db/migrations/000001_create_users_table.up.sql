CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  nickname VARCHAR(20) NOT NULL,
  fullname VARCHAR(20) NOT NULL,
  fb_provider_type VARCHAR(50),
  fb_uid VARCHAR(255),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP
);

ALTER TABLE
  users
ADD
  CONSTRAINT users_email_uix UNIQUE (email);

ALTER TABLE
  users
ADD
  CONSTRAINT users_fb_uid_uix UNIQUE (fb_uid);

CREATE INDEX users_email_idx ON users (email);

CREATE INDEX users_fb_uid_idx ON users (fb_uid);