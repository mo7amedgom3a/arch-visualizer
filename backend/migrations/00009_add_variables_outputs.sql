-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_variables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    project_id UUID NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL, -- string, number, bool, list(...), map(...)
    description TEXT,
    default_value JSONB,
    sensitive BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (project_id, name)
);

CREATE INDEX idx_project_variables_project_id ON project_variables (project_id);

CREATE TABLE project_outputs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    project_id UUID NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    value TEXT NOT NULL, -- The expression
    sensitive BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (project_id, name)
);

CREATE INDEX idx_project_outputs_project_id ON project_outputs (project_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_outputs;

DROP TABLE IF EXISTS project_variables;
-- +goose StatementEnd