-- +goose Up
UPDATE user_stats SET xp = 1 WHERE user_id = '5341aa19-d907-4c01-9454-9d00d1b377ad';

-- +goose Down
