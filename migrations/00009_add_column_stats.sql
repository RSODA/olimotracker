-- +goose Up
ALTER TABLE user_stats
ADD COLUMN is_study_today BOOL DEFAULT FALSE;

-- +goose Down
ALTER TABLE user_stats
DROP COLUMN is_study_today;
