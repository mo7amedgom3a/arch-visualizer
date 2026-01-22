-- +goose Up
-- +goose StatementBegin

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    created_at TIMESTAMP DEFAULT now()
);

-- IaC Targets table (must be created before projects due to FK)
CREATE TABLE iac_targets (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    infra_tool INT REFERENCES iac_targets(id),
    name TEXT NOT NULL,
    cloud_provider TEXT NOT NULL CHECK (cloud_provider IN ('aws', 'azure', 'gcp')),
    region TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

-- Resource Categories table
CREATE TABLE resource_categories (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Resource Kinds table
CREATE TABLE resource_kinds (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Resource Types table
CREATE TABLE resource_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    cloud_provider TEXT NOT NULL,
    category_id INT REFERENCES resource_categories(id),
    kind_id INT REFERENCES resource_kinds(id),
    is_regional BOOLEAN DEFAULT true,
    is_global BOOLEAN DEFAULT false,
    UNIQUE (name, cloud_provider)
);

-- Resources table
CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    resource_type_id INT REFERENCES resource_types(id),
    name TEXT NOT NULL,
    pos_x INT NOT NULL,
    pos_y INT NOT NULL,
    config JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMP DEFAULT now()
);

-- Resource Containment table
CREATE TABLE resource_containment (
    parent_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    child_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    PRIMARY KEY (parent_resource_id, child_resource_id)
);

-- Dependency Types table
CREATE TABLE dependency_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Resource Dependencies table
CREATE TABLE resource_dependencies (
    from_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    to_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    dependency_type_id INT REFERENCES dependency_types(id),
    PRIMARY KEY (from_resource_id, to_resource_id)
);

-- Resource Constraints table
CREATE TABLE resource_constraints (
    id SERIAL PRIMARY KEY,
    resource_type_id INT REFERENCES resource_types(id),
    constraint_type TEXT NOT NULL,
    constraint_value TEXT NOT NULL
);

-- Pricing tables
CREATE TABLE project_pricing (
    id SERIAL PRIMARY KEY,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    total_cost NUMERIC(12, 4) NOT NULL,
    currency TEXT NOT NULL CHECK (currency IN ('USD', 'EUR', 'GBP')),
    period TEXT NOT NULL CHECK (period IN ('hourly', 'monthly', 'yearly')),
    duration_seconds BIGINT NOT NULL,
    provider TEXT NOT NULL CHECK (provider IN ('aws', 'azure', 'gcp')),
    region TEXT,
    calculated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE service_pricing (
    id SERIAL PRIMARY KEY,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    category_id INT REFERENCES resource_categories(id),
    total_cost NUMERIC(12, 4) NOT NULL,
    currency TEXT NOT NULL CHECK (currency IN ('USD', 'EUR', 'GBP')),
    period TEXT NOT NULL CHECK (period IN ('hourly', 'monthly', 'yearly')),
    duration_seconds BIGINT NOT NULL,
    provider TEXT NOT NULL CHECK (provider IN ('aws', 'azure', 'gcp')),
    region TEXT,
    calculated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE service_type_pricing (
    id SERIAL PRIMARY KEY,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    resource_type_id INT REFERENCES resource_types(id),
    total_cost NUMERIC(12, 4) NOT NULL,
    currency TEXT NOT NULL CHECK (currency IN ('USD', 'EUR', 'GBP')),
    period TEXT NOT NULL CHECK (period IN ('hourly', 'monthly', 'yearly')),
    duration_seconds BIGINT NOT NULL,
    provider TEXT NOT NULL CHECK (provider IN ('aws', 'azure', 'gcp')),
    region TEXT,
    calculated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE resource_pricing (
    id SERIAL PRIMARY KEY,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    total_cost NUMERIC(12, 4) NOT NULL,
    currency TEXT NOT NULL CHECK (currency IN ('USD', 'EUR', 'GBP')),
    period TEXT NOT NULL CHECK (period IN ('hourly', 'monthly', 'yearly')),
    duration_seconds BIGINT NOT NULL,
    provider TEXT NOT NULL CHECK (provider IN ('aws', 'azure', 'gcp')),
    region TEXT,
    calculated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE pricing_components (
    id SERIAL PRIMARY KEY,
    resource_pricing_id INT REFERENCES resource_pricing(id) ON DELETE CASCADE,
    component_name TEXT NOT NULL,
    model TEXT NOT NULL CHECK (model IN ('per_hour', 'per_gb', 'per_request', 'one_time', 'tiered', 'percentage')),
    unit TEXT NOT NULL,
    quantity NUMERIC(14, 4) NOT NULL,
    unit_rate NUMERIC(14, 6) NOT NULL,
    subtotal NUMERIC(14, 4) NOT NULL,
    currency TEXT NOT NULL CHECK (currency IN ('USD', 'EUR', 'GBP'))
);

-- Create indexes for better query performance
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_projects_infra_tool ON projects(infra_tool);
CREATE INDEX idx_resources_project_id ON resources(project_id);
CREATE INDEX idx_resources_resource_type_id ON resources(resource_type_id);
CREATE INDEX idx_resource_types_category_id ON resource_types(category_id);
CREATE INDEX idx_resource_types_kind_id ON resource_types(kind_id);
CREATE INDEX idx_resource_constraints_resource_type_id ON resource_constraints(resource_type_id);
CREATE INDEX idx_resource_dependencies_dependency_type_id ON resource_dependencies(dependency_type_id);
CREATE INDEX idx_project_pricing_project_id ON project_pricing(project_id);
CREATE INDEX idx_service_pricing_project_id ON service_pricing(project_id);
CREATE INDEX idx_service_pricing_category_id ON service_pricing(category_id);
CREATE INDEX idx_service_type_pricing_project_id ON service_type_pricing(project_id);
CREATE INDEX idx_service_type_pricing_resource_type_id ON service_type_pricing(resource_type_id);
CREATE INDEX idx_resource_pricing_project_id ON resource_pricing(project_id);
CREATE INDEX idx_resource_pricing_resource_id ON resource_pricing(resource_id);
CREATE INDEX idx_pricing_components_resource_pricing_id ON pricing_components(resource_pricing_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop tables in reverse order to respect foreign key constraints
DROP TABLE IF EXISTS pricing_components;
DROP TABLE IF EXISTS resource_pricing;
DROP TABLE IF EXISTS service_type_pricing;
DROP TABLE IF EXISTS service_pricing;
DROP TABLE IF EXISTS project_pricing;
DROP TABLE IF EXISTS resource_constraints;
DROP TABLE IF EXISTS resource_dependencies;
DROP TABLE IF EXISTS dependency_types;
DROP TABLE IF EXISTS resource_containment;
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS resource_types;
DROP TABLE IF EXISTS resource_kinds;
DROP TABLE IF EXISTS resource_categories;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS iac_targets;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
