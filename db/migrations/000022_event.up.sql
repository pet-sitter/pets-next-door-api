CREATE TABLE IF NOT EXISTS events (
  id UUID PRIMARY KEY,
  event_type VARCHAR NOT NULL,
  author_id UUID NOT NULL,
  name VARCHAR NOT NULL,
  description TEXT NOT NULL,
  media_id UUID,
  topics TEXT[] NOT NULL,
  max_participants INT,
  fee INT NOT NULL,
  start_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMPTZ
);