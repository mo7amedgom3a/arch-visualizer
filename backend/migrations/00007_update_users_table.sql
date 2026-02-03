-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT DEFAULT '';

ALTER TABLE users ADD COLUMN IF NOT EXISTS auth0_id TEXT;

ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN IF EXISTS auth0_id;

ALTER TABLE users DROP COLUMN IF EXISTS avatar;
-- +goose StatementEnd