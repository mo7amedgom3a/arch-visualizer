-- +goose Up
-- +goose StatementBegin

-- Enable UUID extension for marketplace tables
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Update users table to match marketplace schema
ALTER TABLE users ALTER COLUMN id SET DEFAULT uuid_generate_v4();

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS avatar VARCHAR(500),
    ADD COLUMN IF NOT EXISTS is_verified BOOLEAN DEFAULT false,
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;

UPDATE users
SET name = COALESCE(name, 'Unknown');

ALTER TABLE users
    ALTER COLUMN name TYPE VARCHAR(255),
    ALTER COLUMN name SET NOT NULL;

ALTER TABLE users DROP COLUMN IF EXISTS email;

DROP INDEX IF EXISTS idx_users_deleted_at;
ALTER TABLE users DROP COLUMN IF EXISTS deleted_at;

-- Marketplace: Categories table
CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Templates table
CREATE TABLE templates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    cloud_provider VARCHAR(50) NOT NULL CHECK (cloud_provider IN ('AWS', 'Azure', 'GCP', 'Multi-Cloud')),
    rating DECIMAL(3,2) DEFAULT 0 CHECK (rating >= 0 AND rating <= 5),
    review_count INTEGER DEFAULT 0,
    downloads INTEGER DEFAULT 0,
    price DECIMAL(10,2) DEFAULT 0,
    is_subscription BOOLEAN DEFAULT false,
    subscription_price DECIMAL(10,2),
    estimated_cost_min DECIMAL(10,2) NOT NULL,
    estimated_cost_max DECIMAL(10,2) NOT NULL,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Template Technologies junction table (Many-to-Many)
CREATE TABLE template_technologies (
    template_id UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    technology_id UUID NOT NULL REFERENCES technologies(id) ON DELETE CASCADE,
    PRIMARY KEY (template_id, technology_id)
);

-- Marketplace: IAC Formats lookup table
CREATE TABLE iac_formats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Template IAC Formats junction table (Many-to-Many)
CREATE TABLE template_iac_formats (
    template_id UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    iac_format_id UUID NOT NULL REFERENCES iac_formats(id) ON DELETE CASCADE,
    PRIMARY KEY (template_id, iac_format_id)
);

-- Marketplace: Compliance standards lookup table
CREATE TABLE compliance_standards (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Template Compliance junction table (Many-to-Many)
CREATE TABLE template_compliance (
    template_id UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    compliance_id UUID NOT NULL REFERENCES compliance_standards(id) ON DELETE CASCADE,
    PRIMARY KEY (template_id, compliance_id)
);

-- Marketplace: Template Use Cases
CREATE TABLE template_use_cases (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_id UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    icon VARCHAR(100),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Template What You Get items
CREATE TABLE template_features (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_id UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    feature TEXT NOT NULL,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Template Components
CREATE TABLE template_components (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_id UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    service VARCHAR(255) NOT NULL,
    configuration TEXT,
    monthly_cost DECIMAL(10,2) DEFAULT 0,
    purpose TEXT,
    display_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Marketplace: Reviews table
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    template_id UUID NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
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
CREATE INDEX idx_templates_category ON templates(category_id);
CREATE INDEX idx_templates_author ON templates(author_id);
CREATE INDEX idx_templates_cloud_provider ON templates(cloud_provider);
CREATE INDEX idx_templates_rating ON templates(rating DESC);
CREATE INDEX idx_templates_downloads ON templates(downloads DESC);
CREATE INDEX idx_templates_price ON templates(price);
CREATE INDEX idx_templates_is_popular ON templates(is_popular);
CREATE INDEX idx_templates_is_new ON templates(is_new);
CREATE INDEX idx_templates_created_at ON templates(created_at DESC);

CREATE INDEX idx_reviews_template ON reviews(template_id);
CREATE INDEX idx_reviews_user ON reviews(user_id);
CREATE INDEX idx_reviews_rating ON reviews(rating);
CREATE INDEX idx_reviews_created_at ON reviews(created_at DESC);

CREATE INDEX idx_template_technologies_template ON template_technologies(template_id);
CREATE INDEX idx_template_technologies_tech ON template_technologies(technology_id);

CREATE INDEX idx_template_iac_formats_template ON template_iac_formats(template_id);
CREATE INDEX idx_template_compliance_template ON template_compliance(template_id);

CREATE INDEX idx_template_use_cases_template ON template_use_cases(template_id);
CREATE INDEX idx_template_features_template ON template_features(template_id);
CREATE INDEX idx_template_components_template ON template_components(template_id);

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

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop triggers and functions
DROP TRIGGER IF EXISTS update_reviews_updated_at ON reviews;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;
DROP TRIGGER IF EXISTS update_templates_updated_at ON templates;
DROP TRIGGER IF EXISTS trigger_update_template_rating ON reviews;

DROP FUNCTION IF EXISTS update_updated_at_column();
DROP FUNCTION IF EXISTS update_template_rating();

-- Drop marketplace indexes
DROP INDEX IF EXISTS idx_template_components_template;
DROP INDEX IF EXISTS idx_template_features_template;
DROP INDEX IF EXISTS idx_template_use_cases_template;
DROP INDEX IF EXISTS idx_template_compliance_template;
DROP INDEX IF EXISTS idx_template_iac_formats_template;
DROP INDEX IF EXISTS idx_template_technologies_tech;
DROP INDEX IF EXISTS idx_template_technologies_template;
DROP INDEX IF EXISTS idx_reviews_created_at;
DROP INDEX IF EXISTS idx_reviews_rating;
DROP INDEX IF EXISTS idx_reviews_user;
DROP INDEX IF EXISTS idx_reviews_template;
DROP INDEX IF EXISTS idx_templates_created_at;
DROP INDEX IF EXISTS idx_templates_is_new;
DROP INDEX IF EXISTS idx_templates_is_popular;
DROP INDEX IF EXISTS idx_templates_price;
DROP INDEX IF EXISTS idx_templates_downloads;
DROP INDEX IF EXISTS idx_templates_rating;
DROP INDEX IF EXISTS idx_templates_cloud_provider;
DROP INDEX IF EXISTS idx_templates_author;
DROP INDEX IF EXISTS idx_templates_category;

-- Drop marketplace tables (reverse order)
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS template_components;
DROP TABLE IF EXISTS template_features;
DROP TABLE IF EXISTS template_use_cases;
DROP TABLE IF EXISTS template_compliance;
DROP TABLE IF EXISTS template_iac_formats;
DROP TABLE IF EXISTS template_technologies;
DROP TABLE IF EXISTS compliance_standards;
DROP TABLE IF EXISTS iac_formats;
DROP TABLE IF EXISTS technologies;
DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS categories;

-- Revert users table to pre-marketplace schema
ALTER TABLE users ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE users
    ADD COLUMN IF NOT EXISTS email TEXT;

UPDATE users
SET email = COALESCE(email, CONCAT('user+', id, '@example.com'));

ALTER TABLE users
    ALTER COLUMN email SET NOT NULL,
    ALTER COLUMN name TYPE TEXT,
    ALTER COLUMN name DROP NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON users(email);

ALTER TABLE users
    DROP COLUMN IF EXISTS avatar,
    DROP COLUMN IF EXISTS is_verified,
    DROP COLUMN IF EXISTS updated_at;

ALTER TABLE users ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- +goose StatementEnd
