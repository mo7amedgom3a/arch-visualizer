-- +goose Up
-- +goose StatementBegin
-- Add original_id column to resources table to store the frontend node ID
-- This is needed for resolving references in Terraform code generation
ALTER TABLE resources
ADD COLUMN IF NOT EXISTS original_id VARCHAR(255);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_resources_original_id ON resources (original_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_resources_original_id;

ALTER TABLE resources DROP COLUMN IF EXISTS original_id;
-- +goose StatementEnd