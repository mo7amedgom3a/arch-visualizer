-- +goose Up
-- +goose StatementBegin

ALTER TABLE project_versions ADD COLUMN IF NOT EXISTS message TEXT;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE project_versions DROP COLUMN IF EXISTS message;

-- +goose StatementEnd