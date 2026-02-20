package repository_test

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newTestDB creates an in-memory SQLite database for repository tests and
// installs a minimal schema that matches the tables and columns used by the
// repository layer. It deliberately avoids Postgres-specific types and
// defaults so it can run against SQLite.
func newTestDB(t *testing.T, modelsToMigrate ...interface{}) *gorm.DB {
	t.Helper()

	// Use a private in-memory database per test to avoid cross-test interference.
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite db: %v", err)
	}

	// Install a lightweight schema covering all tables that repository tests touch.
	stmts := []string{
		// Core entities
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			name TEXT,
			email TEXT,
			auth0_id TEXT,
			avatar TEXT,
			is_verified INTEGER,
			created_at DATETIME,
			updated_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS projects (
			id TEXT PRIMARY KEY,
			root_project_id TEXT,
			user_id TEXT,
			infra_tool INTEGER,
			name TEXT,
			description TEXT,
			cloud_provider TEXT,
			region TEXT,
			thumbnail TEXT,
			tags TEXT,
			created_at DATETIME,
			updated_at DATETIME,
			deleted_at DATETIME
		);`,

		// Project versions chain (immutable versioning)
		`CREATE TABLE IF NOT EXISTS project_versions (
			id TEXT PRIMARY KEY,
			project_id TEXT,
			parent_version_id TEXT,
			version_number INTEGER NOT NULL DEFAULT 1,
			created_at DATETIME,
			created_by TEXT
		);`,

		// Lookup / metadata
		`CREATE TABLE IF NOT EXISTS iac_targets (
			id INTEGER PRIMARY KEY,
			name TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS resource_categories (
			id INTEGER PRIMARY KEY,
			name TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS resource_kinds (
			id INTEGER PRIMARY KEY,
			name TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS dependency_types (
			id INTEGER PRIMARY KEY,
			name TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS resource_types (
			id INTEGER PRIMARY KEY,
			name TEXT,
			cloud_provider TEXT,
			category_id INTEGER,
			kind_id INTEGER,
			is_regional INTEGER,
			is_global INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS resource_constraints (
			id INTEGER PRIMARY KEY,
			resource_type_id INTEGER,
			constraint_type TEXT,
			constraint_value TEXT
		);`,

		// Resources & relationships
		`CREATE TABLE IF NOT EXISTS resources (
			id TEXT PRIMARY KEY,
			original_id TEXT,
			project_id TEXT,
			resource_type_id INTEGER,
			name TEXT,
			pos_x INTEGER,
			pos_y INTEGER,
			is_visual_only INTEGER,
			config TEXT,
			created_at DATETIME,
			deleted_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS resource_containment (
			parent_resource_id TEXT,
			child_resource_id TEXT,
			PRIMARY KEY (parent_resource_id, child_resource_id)
		);`,
		`CREATE TABLE IF NOT EXISTS resource_dependencies (
			from_resource_id TEXT,
			to_resource_id TEXT,
			dependency_type_id INTEGER,
			PRIMARY KEY (from_resource_id, to_resource_id)
		);`,
		`CREATE TABLE IF NOT EXISTS resource_ui_states (
			resource_id TEXT PRIMARY KEY,
			selected INTEGER,
			expanded INTEGER
		);`,

		// Marketplace
		`CREATE TABLE IF NOT EXISTS categories (
			id TEXT PRIMARY KEY,
			name TEXT,
			slug TEXT,
			created_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS templates (
			id TEXT PRIMARY KEY,
			title TEXT,
			description TEXT,
			category_id TEXT,
			cloud_provider TEXT,
			rating REAL,
			review_count INTEGER,
			downloads INTEGER,
			price REAL,
			is_subscription INTEGER,
			subscription_price REAL,
			estimated_cost_min REAL,
			estimated_cost_max REAL,
			author_id TEXT,
			image_url TEXT,
			is_popular INTEGER,
			is_new INTEGER,
			last_updated DATETIME,
			resources INTEGER,
			deployment_time TEXT,
			regions TEXT,
			created_at DATETIME,
			updated_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS reviews (
			id TEXT PRIMARY KEY,
			template_id TEXT,
			user_id TEXT,
			rating INTEGER,
			title TEXT,
			content TEXT,
			use_case TEXT,
			team_size TEXT,
			deployment_time TEXT,
			helpful_count INTEGER,
			creator_response TEXT,
			created_at DATETIME,
			updated_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS iac_formats (
			id TEXT PRIMARY KEY,
			name TEXT,
			slug TEXT,
			created_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS technologies (
			id TEXT PRIMARY KEY,
			name TEXT,
			slug TEXT,
			created_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS compliance_standards (
			id TEXT PRIMARY KEY,
			name TEXT,
			slug TEXT,
			created_at DATETIME
		);`,

		// Marketplace join tables
		`CREATE TABLE IF NOT EXISTS template_compliance (
			template_id TEXT,
			compliance_id TEXT,
			PRIMARY KEY (template_id, compliance_id)
		);`,
		`CREATE TABLE IF NOT EXISTS template_iac_formats (
			template_id TEXT,
			iac_format_id TEXT,
			PRIMARY KEY (template_id, iac_format_id)
		);`,
		`CREATE TABLE IF NOT EXISTS template_technologies (
			template_id TEXT,
			technology_id TEXT,
			PRIMARY KEY (template_id, technology_id)
		);`,

		// Marketplace extra entities used in tests
		`CREATE TABLE IF NOT EXISTS template_use_cases (
			id TEXT PRIMARY KEY,
			template_id TEXT,
			icon TEXT,
			title TEXT,
			description TEXT,
			display_order INTEGER,
			created_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS template_features (
			id TEXT PRIMARY KEY,
			template_id TEXT,
			feature TEXT,
			display_order INTEGER,
			created_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS template_components (
			id TEXT PRIMARY KEY,
			template_id TEXT,
			name TEXT,
			service TEXT,
			configuration TEXT,
			monthly_cost REAL,
			purpose TEXT,
			display_order INTEGER,
			created_at DATETIME
		);`,

		// Pricing
		`CREATE TABLE IF NOT EXISTS project_pricing (
			id INTEGER PRIMARY KEY,
			project_id TEXT,
			total_cost REAL,
			currency TEXT,
			period TEXT,
			duration_seconds INTEGER,
			provider TEXT,
			region TEXT,
			calculated_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS service_pricing (
			id INTEGER PRIMARY KEY,
			project_id TEXT,
			category_id INTEGER,
			total_cost REAL,
			currency TEXT,
			period TEXT,
			duration_seconds INTEGER,
			provider TEXT,
			region TEXT,
			calculated_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS service_type_pricing (
			id INTEGER PRIMARY KEY,
			project_id TEXT,
			resource_type_id INTEGER,
			total_cost REAL,
			currency TEXT,
			period TEXT,
			duration_seconds INTEGER,
			provider TEXT,
			region TEXT,
			calculated_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS resource_pricing (
			id INTEGER PRIMARY KEY,
			project_id TEXT,
			resource_id TEXT,
			total_cost REAL,
			currency TEXT,
			period TEXT,
			duration_seconds INTEGER,
			provider TEXT,
			region TEXT,
			calculated_at DATETIME
		);`,
		`CREATE TABLE IF NOT EXISTS pricing_components (
			id INTEGER PRIMARY KEY,
			resource_pricing_id INTEGER,
			component_name TEXT,
			model TEXT,
			unit TEXT,
			quantity REAL,
			unit_rate REAL,
			subtotal REAL,
			currency TEXT
		);`,
	}

	for _, stmt := range stmts {
		if err := db.Exec(stmt).Error; err != nil {
			t.Fatalf("failed to install test schema: %v", err)
		}
	}

	return db
}
