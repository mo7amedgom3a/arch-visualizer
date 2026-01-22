-- +goose Up
-- +goose StatementBegin

-- Add deleted_at column to users table for soft deletes
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMP;

-- Add deleted_at column to projects table for soft deletes
ALTER TABLE projects ADD COLUMN deleted_at TIMESTAMP;

-- Add deleted_at column to resources table for soft deletes
ALTER TABLE resources ADD COLUMN deleted_at TIMESTAMP;

-- Create indexes for deleted_at columns to improve query performance
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_projects_deleted_at ON projects(deleted_at);
CREATE INDEX idx_resources_deleted_at ON resources(deleted_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes
DROP INDEX IF EXISTS idx_resources_deleted_at;
DROP INDEX IF EXISTS idx_projects_deleted_at;
DROP INDEX IF EXISTS idx_users_deleted_at;

-- Drop deleted_at columns
ALTER TABLE resources DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE projects DROP COLUMN IF EXISTS deleted_at;
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;

-- +goose StatementEnd
