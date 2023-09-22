CREATE TABLE IF NOT EXISTS breeds
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100) NOT NULL,
    pet_type   VARCHAR(10) NOT NULL,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

ALTER TABLE breeds
    ADD CONSTRAINT breeds_name_pet_type_uix UNIQUE (name, pet_type);
