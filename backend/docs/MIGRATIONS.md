# Database Migrations Guide

This project uses [Goose](https://github.com/pressly/goose) for database migrations.

## Migration File Naming

Migrations are numbered sequentially:
- `00001_init_schema.sql`
- `00002_add_soft_delete.sql`
- `00003_add_marketplace.sql`
- `00004_add_is_visual_only_to_resources.sql`
- etc.

## Migration File Structure

Each migration file must have two sections:

### Up Migration (`-- +goose Up`)
Contains the SQL to apply the migration.

### Down Migration (`-- +goose Down`)
Contains the SQL to rollback the migration.

Example:
```sql
-- +goose Up
-- +goose StatementBegin

-- Your migration SQL here
ALTER TABLE resources ADD COLUMN is_visual_only BOOLEAN DEFAULT false;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Your rollback SQL here
ALTER TABLE resources DROP COLUMN IF EXISTS is_visual_only;

-- +goose StatementEnd
```

## Running Migrations

### Using the platform command:
```bash
go run cmd/platform/main.go
```

### Programmatically:
```go
import "github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"

err := database.RunMigrations("migrations")
```

### Using Goose CLI directly:
```bash
# Install goose if not already installed
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir migrations postgres "host=localhost user=postgres dbname=arch_visualizer sslmode=disable" up

# Rollback last migration
goose -dir migrations postgres "host=localhost user=postgres dbname=arch_visualizer sslmode=disable" down

# Check migration status
goose -dir migrations postgres "host=localhost user=postgres dbname=arch_visualizer sslmode=disable" status
```

## Creating a New Migration

### Step 1: Create the migration file

Create a new file in `migrations/` with the next sequential number:
```bash
touch migrations/00004_add_your_feature.sql
```

### Step 2: Write the Up migration

Add your SQL changes:
```sql
-- +goose Up
-- +goose StatementBegin

-- Example: Add a new column
ALTER TABLE resources ADD COLUMN IF NOT EXISTS new_column TEXT;

-- Example: Create an index
CREATE INDEX IF NOT EXISTS idx_resources_new_column ON resources(new_column);

-- +goose StatementEnd
```

### Step 3: Write the Down migration

Add the rollback SQL:
```sql
-- +goose Down
-- +goose StatementBegin

-- Rollback: Drop index
DROP INDEX IF EXISTS idx_resources_new_column;

-- Rollback: Drop column
ALTER TABLE resources DROP COLUMN IF EXISTS new_column;

-- +goose StatementEnd
```

### Step 4: Test the migration

```bash
# Run migrations
go run cmd/platform/main.go

# Or use goose directly
goose -dir migrations postgres "your-connection-string" up
```

## Common Migration Patterns

### Adding a Column
```sql
-- +goose Up
ALTER TABLE table_name ADD COLUMN IF NOT EXISTS column_name TYPE DEFAULT value;
CREATE INDEX IF NOT EXISTS idx_table_name_column_name ON table_name(column_name);

-- +goose Down
DROP INDEX IF EXISTS idx_table_name_column_name;
ALTER TABLE table_name DROP COLUMN IF EXISTS column_name;
```

### Adding a Foreign Key
```sql
-- +goose Up
ALTER TABLE child_table 
ADD CONSTRAINT fk_child_parent 
FOREIGN KEY (parent_id) REFERENCES parent_table(id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE child_table DROP CONSTRAINT IF EXISTS fk_child_parent;
```

### Creating a Table
```sql
-- +goose Up
CREATE TABLE new_table (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE INDEX idx_new_table_name ON new_table(name);

-- +goose Down
DROP TABLE IF EXISTS new_table;
```

### Modifying a Column
```sql
-- +goose Up
ALTER TABLE table_name ALTER COLUMN column_name TYPE NEW_TYPE;
ALTER TABLE table_name ALTER COLUMN column_name SET NOT NULL;

-- +goose Down
ALTER TABLE table_name ALTER COLUMN column_name DROP NOT NULL;
ALTER TABLE table_name ALTER COLUMN column_name TYPE OLD_TYPE;
```

## Best Practices

1. **Always use `IF EXISTS` / `IF NOT EXISTS`** to make migrations idempotent
2. **Always provide a Down migration** for rollback capability
3. **Test migrations** on a development database first
4. **Use transactions** for complex migrations (Goose handles this automatically)
5. **Keep migrations small** - one logical change per migration
6. **Never modify existing migrations** - create a new one instead
7. **Use descriptive names** that explain what the migration does

## Troubleshooting

### Migration fails with "relation already exists"
- Check if the migration was partially applied
- Use `IF NOT EXISTS` clauses
- Check migration status: `goose status`

### Need to reset migrations
```sql
-- WARNING: This deletes all data!
DROP SCHEMA public CASCADE;
CREATE SCHEMA public;
```
Then re-run all migrations.

### Check migration status
```bash
goose -dir migrations postgres "connection-string" status
```
