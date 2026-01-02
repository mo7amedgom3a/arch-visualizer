CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    infra_tool SERIAL  REFERENCES iac_targets(id),
    name TEXT NOT NULL,
    cloud_provider TEXT NOT NULL CHECK (cloud_provider IN ('aws', 'azure', 'gcp')),
    region TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);


CREATE TABLE resource_categories (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE resource_kinds (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

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

CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    resource_type_id INT REFERENCES resource_types(id),
    name TEXT NOT NULL,

    -- Visual positioning
    pos_x INT NOT NULL,
    pos_y INT NOT NULL,

    -- JSON config (CIDR, instance_type, tags...)
    config JSONB NOT NULL DEFAULT '{}',

    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE resource_containment (
    parent_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    child_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    PRIMARY KEY (parent_resource_id, child_resource_id)
);

CREATE TABLE dependency_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE resource_dependencies (
    from_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    to_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    dependency_type_id INT REFERENCES dependency_types(id),
    PRIMARY KEY (from_resource_id, to_resource_id)
);

CREATE TABLE resource_constraints (
    id SERIAL PRIMARY KEY,
    resource_type_id INT REFERENCES resource_types(id),
    constraint_type TEXT NOT NULL,
    constraint_value TEXT NOT NULL
);

CREATE TABLE iac_targets (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);
