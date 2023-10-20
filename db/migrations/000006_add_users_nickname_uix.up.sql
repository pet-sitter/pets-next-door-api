ALTER TABLE
    users
    ADD
        CONSTRAINT users_nickname_uix UNIQUE (nickname);
