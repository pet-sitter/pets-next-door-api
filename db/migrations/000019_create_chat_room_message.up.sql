CREATE TABLE rooms (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    room_type VARCHAR(255) NOT NULL,
    group_id BIGINT NOT NULL,
    created_at TYPE timestamptz USING created_at AT TIME ZONE 'UTC',
    updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC',
    deleted_at TYPE timestamptz USING deleted_at AT TIME ZONE 'UTC';
);

CREATE TABLE messages (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    room_id BIGINT NOT NULL,
    message_type VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TYPE timestamptz USING created_at AT TIME ZONE 'UTC',
    updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC',
    deleted_at TYPE timestamptz USING deleted_at AT TIME ZONE 'UTC',
    FOREIGN KEY (RoomID) REFERENCES Rooms(ID);
);

CREATE TABLE user_chat_messages (
    id BIGINT PRIMARY KEY,
    user_id BIGINT NOT NULL,
    room_id BIGINT NOT NULL,
    created_at TYPE timestamptz USING created_at AT TIME ZONE 'UTC',
    updated_at TYPE timestamptz USING updated_at AT TIME ZONE 'UTC',
    deleted_at TYPE timestamptz USING deleted_at AT TIME ZONE 'UTC',
    FOREIGN KEY (RoomID) REFERENCES Rooms(ID);
);
