-- +goose Up
ALTER TABLE user_stats
ADD COLUMN goal INT DEFAULT 60;

-- +goose Down
ALTER TABLE user_stats
DROP COLUMN goal;
