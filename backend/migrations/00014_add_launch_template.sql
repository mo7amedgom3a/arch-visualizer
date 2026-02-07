-- +goose Up
-- +goose StatementBegin

-- Ensure 'Configuration' kind exists
INSERT INTO
    resource_kinds (name)
SELECT 'Configuration'
WHERE
    NOT EXISTS (
        SELECT 1
        FROM resource_kinds
        WHERE
            name = 'Configuration'
    );

-- Get IDs for Category and Kind
DO $$
DECLARE
    compute_category_id INTEGER;
    configuration_kind_id INTEGER;
BEGIN
    SELECT id INTO compute_category_id FROM resource_categories WHERE name = 'Compute';
    SELECT id INTO configuration_kind_id FROM resource_kinds WHERE name = 'Configuration';

    -- Insert LaunchTemplate resource type if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM resource_types WHERE name = 'LaunchTemplate' AND cloud_provider = 'aws') THEN
        INSERT INTO resource_types (name, cloud_provider, category_id, kind_id, is_regional, is_global)
        VALUES ('LaunchTemplate', 'aws', compute_category_id, configuration_kind_id, true, false);
    END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM resource_types
WHERE
    name = 'LaunchTemplate'
    AND cloud_provider = 'aws';
-- +goose StatementEnd