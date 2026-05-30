-- +goose Up
ALTER TABLE users ADD COLUMN username VARCHAR(55) NOT NULL UNIQUE;

-- +goose Down
ALTER TABLE users DROP COLUMN username;
