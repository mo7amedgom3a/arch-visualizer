-- +goose Up
-- +goose StatementBegin

-- Ensure 'Networking' category exists
INSERT INTO
    resource_categories (name)
SELECT 'Networking'
WHERE
    NOT EXISTS (
        SELECT 1
        FROM resource_categories
        WHERE
            name = 'Networking'
    );

-- Ensure 'Zone' kind exists
INSERT INTO
    resource_kinds (name)
SELECT 'Zone'
WHERE
    NOT EXISTS (
        SELECT 1
        FROM resource_kinds
        WHERE
            name = 'Zone'
    );

-- Get IDs for Category and Kind
DO $$
DECLARE
    networking_category_id INTEGER;
    zone_kind_id INTEGER;
BEGIN
    SELECT id INTO networking_category_id FROM resource_categories WHERE name = 'Networking';
    SELECT id INTO zone_kind_id FROM resource_kinds WHERE name = 'Zone';

    -- Insert AvailabilityZone resource type if it doesn't exist
    IF NOT EXISTS (SELECT 1 FROM resource_types WHERE name = 'AvailabilityZone' AND cloud_provider = 'aws') THEN
        INSERT INTO resource_types (name, cloud_provider, category_id, kind_id, is_regional, is_global)
        VALUES ('AvailabilityZone', 'aws', networking_category_id, zone_kind_id, true, false);
    END IF;
END $$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM resource_types
WHERE
    name = 'AvailabilityZone'
    AND cloud_provider = 'aws';
-- +goose StatementEnd