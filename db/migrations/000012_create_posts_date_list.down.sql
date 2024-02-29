ALTER TABLE sos_posts
    ADD COLUMN date_start_at DATE,
    ADD COLUMN date_end_at DATE;

DROP INDEX sos_posts_dates_sos_post_id;
DROP TABLE IF EXISTS sos_posts_dates;
DROP TABLE IF EXISTS sos_dates;