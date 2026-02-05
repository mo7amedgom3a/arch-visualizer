-- +goose Up
-- +goose StatementBegin

-- Resource UI States table
CREATE TABLE resource_ui_states (
    id SERIAL PRIMARY KEY,
    resource_id UUID NOT NULL REFERENCES resources (id) ON DELETE CASCADE,
    x FLOAT NOT NULL DEFAULT 0,
    y FLOAT NOT NULL DEFAULT 0,
    width FLOAT,
    height FLOAT,
    style JSONB,
    measured JSONB,
    selected BOOLEAN DEFAULT false,
    dragging BOOLEAN DEFAULT false,
    resizing BOOLEAN DEFAULT false,
    focusable BOOLEAN DEFAULT true,
    selectable BOOLEAN DEFAULT true,
    z_index INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE (resource_id)
);

CREATE INDEX idx_resource_ui_states_resource_id ON resource_ui_states (resource_id);

-- Project UI States table (for global view state)
CREATE TABLE project_ui_states (
    id SERIAL PRIMARY KEY,
    project_id UUID NOT NULL REFERENCES projects (id) ON DELETE CASCADE,
    zoom FLOAT DEFAULT 1.0,
    viewport_x FLOAT DEFAULT 0,
    viewport_y FLOAT DEFAULT 0,
    selected_node_ids JSONB DEFAULT '[]',
    selected_edge_ids JSONB DEFAULT '[]',
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE (project_id)
);

CREATE INDEX idx_project_ui_states_project_id ON project_ui_states (project_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS project_ui_states;

DROP TABLE IF EXISTS resource_ui_states;

-- +goose StatementEnd