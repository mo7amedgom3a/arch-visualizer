-- +goose Up
-- +goose StatementBegin

-- Add is_visual_only column to resources table
-- This flag tracks whether a resource is visual-only (icon) vs real infrastructure
ALTER TABLE resources ADD COLUMN IF NOT EXISTS is_visual_only BOOLEAN DEFAULT false;

-- Create index for better query performance when filtering visual-only resources
CREATE INDEX IF NOT EXISTS idx_resources_is_visual_only ON resources(is_visual_only);

-- Update existing resources to set is_visual_only based on config metadata
-- If config contains "isVisualOnly": true, set the column to true
UPDATE resources
SET is_visual_only = COALESCE((config->>'isVisualOnly')::boolean, false)
WHERE config ? 'isVisualOnly';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop index
DROP INDEX IF EXISTS idx_resources_is_visual_only;

-- Drop column
ALTER TABLE resources DROP COLUMN IF EXISTS is_visual_only;

-- +goose StatementEnd
