ALTER TABLE
users
ADD COLUMN uuid UUID NULL,
ADD COLUMN profile_image_uuid UUID NULL;

ALTER TABLE
media
ADD COLUMN uuid UUID NULL;

ALTER TABLE
breeds
ADD COLUMN uuid UUID NULL;

ALTER TABLE
pets
ADD COLUMN uuid UUID NULL,
ADD COLUMN owner_uuid UUID NULL,
ADD COLUMN profile_image_uuid UUID NULL;

ALTER TABLE
base_posts
ADD COLUMN uuid UUID NULL,
ADD COLUMN author_uuid UUID NULL;;

ALTER TABLE
sos_posts
ADD COLUMN thumbnail_uuid UUID NULL;

ALTER TABLE
sos_dates
ADD COLUMN uuid UUID NULL;

ALTER TABLE
sos_posts_dates
ADD COLUMN uuid UUID NULL,
ADD COLUMN sos_post_uuid UUID NULL,
ADD COLUMN sos_dates_uuid UUID NULL;

ALTER TABLE
sos_conditions
ADD COLUMN uuid UUID NULL;

ALTER TABLE
sos_posts_conditions
ADD COLUMN uuid UUID NULL,
ADD COLUMN sos_post_uuid UUID NULL,
ADD COLUMN sos_condition_uuid UUID NULL;

ALTER TABLE
sos_posts_pets
ADD COLUMN uuid UUID NULL,
ADD COLUMN sos_post_uuid UUID NULL,
ADD COLUMN pet_uuid UUID NULL;

ALTER TABLE
resource_media
ADD COLUMN uuid UUID NULL,
ADD COLUMN media_uuid UUID NULL,
ADD COLUMN resource_uuid UUID NULL;;

