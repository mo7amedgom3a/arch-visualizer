-- +goose Up
-- +goose StatementBegin
ALTER TABLE resources DROP COLUMN IF EXISTS pos_x;

ALTER TABLE resources DROP COLUMN IF EXISTS pos_y;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE resources ADD COLUMN pos_x INT NOT NULL DEFAULT 0;

ALTER TABLE resources ADD COLUMN pos_y INT NOT NULL DEFAULT 0;
-- +goose StatementEnd