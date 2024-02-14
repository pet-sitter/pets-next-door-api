DROP INDEX users_email_uix_null;

DROP INDEX users_fb_uid_uix_null;

DROP INDEX users_nickname_uix_null;

DROP INDEX users_email_idx;

DROP INDEX users_fb_uid_idx;

ALTER TABLE
  users
ADD
  CONSTRAINT users_email_uix UNIQUE (email),
ADD
  CONSTRAINT users_fb_uid_uix UNIQUE (fb_uid),
ADD
  CONSTRAINT users_nickname_uix UNIQUE (nickname);

CREATE INDEX users_email_idx ON users (email);

CREATE INDEX users_fb_uid_idx ON users (fb_uid);