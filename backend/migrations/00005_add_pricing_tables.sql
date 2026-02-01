-- +goose Up
-- +goose StatementBegin

-- Create pricing_rates table for scalable pricing system
-- This table stores pricing rates per resource type, allowing dynamic pricing updates
CREATE TABLE IF NOT EXISTS pricing_rates (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(20) NOT NULL,
    resource_type VARCHAR(100) NOT NULL,
    component_name VARCHAR(100) NOT NULL,
    pricing_model VARCHAR(50) NOT NULL,
    CHECK (pricing_model IN ('per_hour', 'per_gb', 'per_request', 'one_time', 'tiered', 'percentage')),
    unit VARCHAR(50) NOT NULL,
    rate NUMERIC(14, 6) NOT NULL,
    currency VARCHAR(10) DEFAULT 'USD' CHECK (currency IN ('USD', 'EUR', 'GBP')),
    region VARCHAR(50),
    effective_from TIMESTAMP NOT NULL DEFAULT NOW(),
    effective_to TIMESTAMP,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create unique constraint using partial index for NULL region handling
CREATE UNIQUE INDEX IF NOT EXISTS unique_pricing_rate 
ON pricing_rates(provider, resource_type, component_name, COALESCE(region, ''), effective_from);

-- Create hidden_dependencies table for implicit resource dependencies
-- This table defines hidden costs (e.g., NAT Gateway automatically includes Elastic IP)
CREATE TABLE IF NOT EXISTS hidden_dependencies (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(20) NOT NULL,
    parent_resource_type VARCHAR(100) NOT NULL,
    child_resource_type VARCHAR(100) NOT NULL,
    quantity_expression VARCHAR(255) DEFAULT '1',
    condition_expression VARCHAR(255),
    is_attached BOOLEAN DEFAULT true,
    description TEXT,
    CONSTRAINT unique_hidden_dependency UNIQUE(provider, parent_resource_type, child_resource_type)
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_pricing_rates_provider ON pricing_rates(provider);
CREATE INDEX IF NOT EXISTS idx_pricing_rates_resource_type ON pricing_rates(resource_type);
CREATE INDEX IF NOT EXISTS idx_pricing_rates_region ON pricing_rates(region);
CREATE INDEX IF NOT EXISTS idx_pricing_rates_effective_from ON pricing_rates(effective_from);
CREATE INDEX IF NOT EXISTS idx_pricing_rates_effective_to ON pricing_rates(effective_to);
CREATE INDEX IF NOT EXISTS idx_pricing_rates_provider_resource ON pricing_rates(provider, resource_type);

CREATE INDEX IF NOT EXISTS idx_hidden_dependencies_provider ON hidden_dependencies(provider);
CREATE INDEX IF NOT EXISTS idx_hidden_dependencies_parent ON hidden_dependencies(parent_resource_type);
CREATE INDEX IF NOT EXISTS idx_hidden_dependencies_child ON hidden_dependencies(child_resource_type);
CREATE INDEX IF NOT EXISTS idx_hidden_dependencies_provider_parent ON hidden_dependencies(provider, parent_resource_type);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes
DROP INDEX IF EXISTS idx_hidden_dependencies_provider_parent;
DROP INDEX IF EXISTS idx_hidden_dependencies_child;
DROP INDEX IF EXISTS idx_hidden_dependencies_parent;
DROP INDEX IF EXISTS idx_hidden_dependencies_provider;

DROP INDEX IF EXISTS unique_pricing_rate;
DROP INDEX IF EXISTS idx_pricing_rates_provider_resource;
DROP INDEX IF EXISTS idx_pricing_rates_effective_to;
DROP INDEX IF EXISTS idx_pricing_rates_effective_from;
DROP INDEX IF EXISTS idx_pricing_rates_region;
DROP INDEX IF EXISTS idx_pricing_rates_resource_type;
DROP INDEX IF EXISTS idx_pricing_rates_provider;

-- Drop tables
DROP TABLE IF EXISTS hidden_dependencies;
DROP TABLE IF EXISTS pricing_rates;

-- +goose StatementEnd
