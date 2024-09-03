ALTER TABLE	users DROP COLUMN uuid, DROP COLUMN profile_image_uuid;
ALTER TABLE	media DROP COLUMN uuid;
ALTER TABLE	breeds DROP COLUMN uuid;
ALTER TABLE	pets DROP COLUMN uuid, DROP COLUMN owner_uuid, DROP COLUMN profile_image_uuid;
ALTER TABLE	base_posts DROP COLUMN uuid, DROP COLUMN author_uuid;
ALTER TABLE	sos_posts DROP COLUMN thumbnail_uuid;
ALTER TABLE	sos_dates DROP COLUMN uuid;
ALTER TABLE	sos_posts_dates DROP COLUMN uuid, DROP COLUMN sos_post_uuid, DROP COLUMN sos_dates_uuid;
ALTER TABLE	sos_conditions DROP COLUMN uuid;
ALTER TABLE	sos_posts_conditions DROP COLUMN uuid, DROP COLUMN sos_post_uuid, DROP COLUMN sos_condition_uuid;
ALTER TABLE	sos_posts_pets DROP COLUMN uuid, DROP COLUMN sos_post_uuid, DROP COLUMN pet_uuid;
ALTER TABLE	resource_media DROP COLUMN uuid, DROP COLUMN media_uuid, DROP COLUMN resource_uuid;

