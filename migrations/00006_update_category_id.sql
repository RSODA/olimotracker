-- +goose Up
ALTER TABLE sessions ALTER COLUMN category_id DROP NOT NULL;

-- +goose Down
ALTER TABLE sessions ALTER COLUMN category_id SET NOT NULL;