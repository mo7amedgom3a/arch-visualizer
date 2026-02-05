-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS resource_constraints;

CREATE TABLE resource_constraints (
    id SERIAL PRIMARY KEY,
    resource_type_id INTEGER NOT NULL REFERENCES resource_types (id) ON DELETE CASCADE,
    constraint_type TEXT NOT NULL,
    constraint_value TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_resource_constraints_resource_type_id ON resource_constraints (resource_type_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE resource_constraints;
-- +goose StatementEnd