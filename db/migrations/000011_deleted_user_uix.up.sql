-- delete previous index and unique constraint
ALTER TABLE
  users DROP CONSTRAINT users_email_uix,
  DROP CONSTRAINT users_fb_uid_uix,
  DROP CONSTRAINT users_nickname_uix;

DROP INDEX users_email_idx;

DROP INDEX users_fb_uid_idx;

-- add the index and unique constraint with deleted_at
CREATE UNIQUE INDEX users_email_uix_null ON users (email)
WHERE
  deleted_at IS NULL;

CREATE UNIQUE INDEX users_fb_uid_uix_null ON users (fb_uid)
WHERE
  deleted_at IS NULL;

CREATE UNIQUE INDEX users_nickname_uix_null ON users (nickname)
WHERE
  deleted_at IS NULL;

CREATE INDEX users_email_idx ON users (email);

CREATE INDEX users_fb_uid_idx ON users (fb_uid);