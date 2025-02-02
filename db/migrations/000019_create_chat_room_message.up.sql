CREATE TABLE chat_rooms (
    id uuid PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    host_id uuid NULL, -- 하위호환 유지를 위해 NULL 허용
    room_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE chat_messages (
    id uuid PRIMARY KEY,
    user_id BIGINT NOT NULL,
    room_id BIGINT NOT NULL,
    message_type VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMPTZ,
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id)
);

CREATE TABLE user_chat_rooms (
    id uuid PRIMARY KEY,
    user_id BIGINT NOT NULL,
    room_id BIGINT NOT NULL,
    joined_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    left_at TIMESTAMPTZ,
    FOREIGN KEY (room_id) REFERENCES chat_rooms(id)
);
