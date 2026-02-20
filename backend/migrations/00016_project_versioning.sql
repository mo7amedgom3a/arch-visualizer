-- +goose Up
-- +goose StatementBegin

-- Drop the old project_versions table (old schema had changes/snapshot columns)
DROP TABLE IF EXISTS project_versions;

-- Recreate project_versions with the new immutable versioning schema
-- Each entry links a project snapshot (project_id) to the previous version (parent_version_id)
CREATE TABLE project_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    project_id UUID NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    parent_version_id UUID REFERENCES project_versions (id),
    version_number INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT now(),
    created_by UUID REFERENCES users (id)
);

CREATE INDEX IF NOT EXISTS idx_project_versions_project_id ON project_versions (project_id);

CREATE INDEX IF NOT EXISTS idx_project_versions_parent_version_id ON project_versions (parent_version_id);

-- Add root_project_id to projects for fast lineage queries
-- NULL means this row IS the root of its lineage
ALTER TABLE projects
ADD COLUMN IF NOT EXISTS root_project_id UUID REFERENCES projects (id);

CREATE INDEX IF NOT EXISTS idx_projects_root_project_id ON projects (root_project_id);

-- Backfill: create version_number=1 for every existing project that has no version entry yet
INSERT INTO
    project_versions (
        id,
        project_id,
        parent_version_id,
        version_number,
        created_at,
        created_by
    )
SELECT gen_random_uuid (), p.id, NULL, 1, p.created_at, p.user_id
FROM projects p
WHERE
    NOT EXISTS (
        SELECT 1
        FROM project_versions pv
        WHERE
            pv.project_id = p.id
    );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Remove the new columns/tables
DROP INDEX IF EXISTS idx_projects_root_project_id;

ALTER TABLE projects DROP COLUMN IF EXISTS root_project_id;

DROP TABLE IF EXISTS project_versions;

-- +goose StatementEnd