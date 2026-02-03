-- Enable UUID extension for marketplace tables
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(255) NOT NULL,
    avatar VARCHAR(500),
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
# Projects example: project-1, project-2, project-3
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    infra_tool SERIAL  REFERENCES iac_targets(id),
    name TEXT NOT NULL,
    description TEXT,
    thumbnail TEXT,
    tags TEXT[],
    cloud_provider TEXT NOT NULL CHECK (cloud_provider IN ('aws', 'azure', 'gcp')),
    region TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    deleted_at TIMESTAMP
);
# Project Versions for version control
CREATE TABLE project_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT now(),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    changes TEXT,
    snapshot JSONB
);

CREATE INDEX idx_project_versions_project_id ON project_versions (project_id);

# Project Variables for Terraform/IaC variables
CREATE TABLE project_variables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT NOT NULL, -- string, number, bool, list(string), map(string)
    description TEXT,
    default_value JSONB,
    sensitive BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),
    UNIQUE(project_id, name)
);

CREATE INDEX idx_project_variables_project_id ON project_variables (project_id);

# Resource Categories example: Compute, Networking, Storage, Database, Security
CREATE TABLE resource_categories (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

# Resource Kinds example: VirtualMachine, Container, Function, Network, LoadBalancer
CREATE TABLE resource_kinds (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

# Resource Types example: EC2, Lambda, S3, RDS, VPC
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

# Resources example: subnet-1, vpc-1, ec2-1, lambda-1, s3-1, rds-1
CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    resource_type_id INT REFERENCES resource_types(id),
    name TEXT NOT NULL,

-- Visual positioning
pos_x INT NOT NULL, pos_y INT NOT NULL,

-- JSON config (CIDR, instance_type, tags...)


config JSONB NOT NULL DEFAULT '{}',

    created_at TIMESTAMP DEFAULT now()
);

# Resource Containment example: VPC → Subnet → EC2
CREATE TABLE resource_containment (
    parent_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    child_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    PRIMARY KEY (parent_resource_id, child_resource_id)
);

# Dependency Types example: uses, depends_on, connects_to, references
CREATE TABLE dependency_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

# Resource Dependencies example: subnet-1 → vpc-1
CREATE TABLE resource_dependencies (
    from_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    to_resource_id UUID REFERENCES resources(id) ON DELETE CASCADE,
    dependency_type_id INT REFERENCES dependency_types(id),
    PRIMARY KEY (from_resource_id, to_resource_id)
);

# Resource Constraints example: subnet must be inside a vpc
CREATE TABLE resource_constraints (
    id SERIAL PRIMARY KEY,
    resource_type_id INT REFERENCES resource_types(id),
    constraint_type TEXT NOT NULL,
    constraint_value TEXT NOT NULL
);

# IaC Targets example: Terraform, Pulumi, CDK
CREATE TABLE iac_targets (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

# Pricing estimates for projects and services/types
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

# Pricing per service (resource category)
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

# Pricing per service type (resource type)
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

# Pricing per resource
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

# Pricing breakdown per component (e.g., per-hour, per-GB, per-request)
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

# Pricing rates table - stores pricing rates per resource type for scalable pricing
CREATE TABLE pricing_rates (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(20) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    component_name VARCHAR(100) NOT NULL,
    pricing_model VARCHAR(50) NOT NULL,
    unit VARCHAR(50) NOT NULL,
    rate NUMERIC(14, 6) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD',
    region VARCHAR(50),
    effective_from TIMESTAMP NOT NULL DEFAULT NOW(),
    effective_to TIMESTAMP,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create unique constraint using partial index for NULL region handling
CREATE UNIQUE INDEX unique_pricing_rate ON pricing_rates (
    provider,
    resource_type,
    component_name,
    COALESCE(region, ''),
    effective_from
);

# Hidden dependencies - defines implicit resource dependencies (e.g., NAT Gateway -> Elastic IP)
CREATE TABLE hidden_dependencies (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(20) NOT NULL,
    parent_resource_type VARCHAR(100) NOT NULL,
    child_resource_type VARCHAR(100) NOT NULL,
    quantity_expression VARCHAR(255) DEFAULT '1',
    condition_expression VARCHAR(255),
    is_attached BOOLEAN DEFAULT true,
    description TEXT,
    UNIQUE(provider, parent_resource_type, child_resource_type)
);

-- Marketplace: Categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Templates table
CREATE TABLE templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category_id UUID NOT NULL REFERENCES categories (id) ON DELETE RESTRICT,
    cloud_provider VARCHAR(50) NOT NULL CHECK (
        cloud_provider IN (
            'AWS',
            'Azure',
            'GCP',
            'Multi-Cloud'
        )
    ),
    rating DECIMAL(3, 2) DEFAULT 0 CHECK (
        rating >= 0
        AND rating <= 5
    ),
    review_count INTEGER DEFAULT 0,
    downloads INTEGER DEFAULT 0,
    price DECIMAL(10, 2) DEFAULT 0,
    is_subscription BOOLEAN DEFAULT false,
    subscription_price DECIMAL(10, 2),
    estimated_cost_min DECIMAL(10, 2) NOT NULL,
    estimated_cost_max DECIMAL(10, 2) NOT NULL,
    author_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    image_url VARCHAR(500),
    is_popular BOOLEAN DEFAULT false,
    is_new BOOLEAN DEFAULT false,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    resources INTEGER DEFAULT 0,
    deployment_time VARCHAR(50),
    regions TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Technologies lookup table
CREATE TABLE technologies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Template Technologies junction table (Many-to-Many)
CREATE TABLE template_technologies (
    template_id UUID NOT NULL REFERENCES templates (id) ON DELETE CASCADE,
    technology_id UUID NOT NULL REFERENCES technologies (id) ON DELETE CASCADE,
    PRIMARY KEY (template_id, technology_id)
);

-- Marketplace: IAC Formats lookup table
CREATE TABLE iac_formats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Template IAC Formats junction table (Many-to-Many)
CREATE TABLE template_iac_formats (
    template_id UUID NOT NULL REFERENCES templates (id) ON DELETE CASCADE,
    iac_format_id UUID NOT NULL REFERENCES iac_formats (id) ON DELETE CASCADE,
    PRIMARY KEY (template_id, iac_format_id)
);

-- Marketplace: Compliance standards lookup table
CREATE TABLE compliance_standards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Template Compliance junction table (Many-to-Many)
CREATE TABLE template_compliance (
    template_id UUID NOT NULL REFERENCES templates (id) ON DELETE CASCADE,
    compliance_id UUID NOT NULL REFERENCES compliance_standards (id) ON DELETE CASCADE,
    PRIMARY KEY (template_id, compliance_id)
);

-- Marketplace: Template Use Cases
CREATE TABLE template_use_cases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    template_id UUID NOT NULL REFERENCES templates (id) ON DELETE CASCADE,
    icon VARCHAR(100),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Template What You Get items
CREATE TABLE template_features (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    template_id UUID NOT NULL REFERENCES templates (id) ON DELETE CASCADE,
    feature TEXT NOT NULL,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Template Components
CREATE TABLE template_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    template_id UUID NOT NULL REFERENCES templates (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    service VARCHAR(255) NOT NULL,
    configuration TEXT,
    monthly_cost DECIMAL(10, 2) DEFAULT 0,
    purpose TEXT,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Reviews table
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    template_id UUID NOT NULL REFERENCES templates (id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (
        rating >= 1
        AND rating <= 5
    ),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    use_case VARCHAR(255),
    team_size VARCHAR(50),
    deployment_time VARCHAR(50),
    helpful_count INTEGER DEFAULT 0,
    creator_response TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace indexes
CREATE INDEX idx_templates_category ON templates (category_id);

CREATE INDEX idx_templates_author ON templates (author_id);

CREATE INDEX idx_templates_cloud_provider ON templates (cloud_provider);

CREATE INDEX idx_templates_rating ON templates (rating DESC);

CREATE INDEX idx_templates_downloads ON templates (downloads DESC);

CREATE INDEX idx_templates_price ON templates (price);

CREATE INDEX idx_templates_is_popular ON templates (is_popular);

CREATE INDEX idx_templates_is_new ON templates (is_new);

CREATE INDEX idx_templates_created_at ON templates (created_at DESC);

CREATE INDEX idx_reviews_template ON reviews (template_id);

CREATE INDEX idx_reviews_user ON reviews (user_id);

CREATE INDEX idx_reviews_rating ON reviews (rating);

CREATE INDEX idx_reviews_created_at ON reviews (created_at DESC);

CREATE INDEX idx_template_technologies_template ON template_technologies (template_id);

CREATE INDEX idx_template_technologies_tech ON template_technologies (technology_id);

CREATE INDEX idx_template_iac_formats_template ON template_iac_formats (template_id);

CREATE INDEX idx_template_compliance_template ON template_compliance (template_id);

CREATE INDEX idx_template_use_cases_template ON template_use_cases (template_id);

CREATE INDEX idx_template_features_template ON template_features (template_id);

CREATE INDEX idx_template_components_template ON template_components (template_id);

-- Pricing tables indexes
CREATE INDEX idx_pricing_rates_provider ON pricing_rates (provider);

CREATE INDEX idx_pricing_rates_resource_type ON pricing_rates (resource_type);

CREATE INDEX idx_pricing_rates_region ON pricing_rates (region);

CREATE INDEX idx_pricing_rates_effective_from ON pricing_rates (effective_from);

CREATE INDEX idx_pricing_rates_effective_to ON pricing_rates (effective_to);

CREATE INDEX idx_hidden_dependencies_provider ON hidden_dependencies (provider);

CREATE INDEX idx_hidden_dependencies_parent ON hidden_dependencies (parent_resource_type);

CREATE INDEX idx_hidden_dependencies_child ON hidden_dependencies (child_resource_type);

-- Marketplace: triggers to update derived fields
CREATE OR REPLACE FUNCTION update_template_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE templates
    SET
        rating = (
            SELECT COALESCE(AVG(rating), 0)
            FROM reviews
            WHERE template_id = COALESCE(NEW.template_id, OLD.template_id)
        ),
        review_count = (
            SELECT COUNT(*)
            FROM reviews
            WHERE template_id = COALESCE(NEW.template_id, OLD.template_id)
        ),
        updated_at = CURRENT_TIMESTAMP
    WHERE id = COALESCE(NEW.template_id, OLD.template_id);

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_template_rating
AFTER INSERT OR UPDATE OR DELETE ON reviews
FOR EACH ROW
EXECUTE FUNCTION update_template_rating();

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_templates_updated_at
BEFORE UPDATE ON templates
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reviews_updated_at
BEFORE UPDATE ON reviews
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();