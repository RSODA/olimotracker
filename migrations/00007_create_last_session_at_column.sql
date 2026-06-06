-- +goose Up
ALTER TABLE user_stats
ADD COLUMN last_session_at TIMESTAMPTZ;

-- +goose Down
ALTER TABLE user_stats
DROP COLUMN last_session_at;
