-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects
ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT now();

ALTER TABLE projects ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects DROP COLUMN IF EXISTS updated_at;

ALTER TABLE projects DROP COLUMN IF EXISTS deleted_at;
-- +goose StatementEnd