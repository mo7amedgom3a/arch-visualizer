-- +goose Up
-- +goose StatementBegin

-- Add instance_type and operating_system columns to pricing_rates table
-- These fields enable instance-type-specific pricing lookups for EC2
ALTER TABLE pricing_rates 
ADD COLUMN IF NOT EXISTS instance_type VARCHAR(50),
ADD COLUMN IF NOT EXISTS operating_system VARCHAR(20) DEFAULT 'linux';

-- Create index for EC2 instance type lookups
CREATE INDEX IF NOT EXISTS idx_pricing_rates_instance_type ON pricing_rates(instance_type);
CREATE INDEX IF NOT EXISTS idx_pricing_rates_operating_system ON pricing_rates(operating_system);

-- Create composite index for efficient EC2 pricing lookups
-- This index supports queries like: FindByInstanceType(provider, instanceType, region, os)
CREATE INDEX IF NOT EXISTS idx_pricing_rates_ec2_lookup 
ON pricing_rates(provider, resource_type, instance_type, region, operating_system)
WHERE instance_type IS NOT NULL;

-- Update the unique constraint to include instance_type and operating_system
-- Drop the old unique index
DROP INDEX IF EXISTS unique_pricing_rate;

-- Create new unique index that includes instance_type and operating_system
CREATE UNIQUE INDEX IF NOT EXISTS unique_pricing_rate 
ON pricing_rates(
    provider, 
    resource_type, 
    component_name, 
    COALESCE(region, ''), 
    COALESCE(instance_type, ''),
    COALESCE(operating_system, 'linux'),
    effective_from
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop indexes
DROP INDEX IF EXISTS unique_pricing_rate;
DROP INDEX IF EXISTS idx_pricing_rates_ec2_lookup;
DROP INDEX IF EXISTS idx_pricing_rates_operating_system;
DROP INDEX IF EXISTS idx_pricing_rates_instance_type;

-- Remove columns
ALTER TABLE pricing_rates 
DROP COLUMN IF EXISTS operating_system,
DROP COLUMN IF EXISTS instance_type;

-- Recreate the original unique index
CREATE UNIQUE INDEX IF NOT EXISTS unique_pricing_rate 
ON pricing_rates(provider, resource_type, component_name, COALESCE(region, ''), effective_from);

-- +goose StatementEnd
