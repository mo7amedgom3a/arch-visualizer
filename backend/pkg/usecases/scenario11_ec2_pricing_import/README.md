# Scenario 11: EC2 Pricing Import and Calculation Test

This scenario tests the new EC2 pricing import feature that integrates scraper data into the pricing system.

## Overview

This usecase demonstrates:
1. **Database Migration**: Runs migration to add `instance_type` and `operating_system` fields to `pricing_rates` table
2. **Pricing Import**: Imports EC2 On-Demand pricing data from scraper JSON output
3. **Pricing Calculation**: Tests pricing calculation with different EC2 instance types, regions, and operating systems
4. **Comparison**: Shows the difference between DB rates (from scraper) and hardcoded fallback rates
5. **Architecture Test**: Creates a test architecture with multiple EC2 instance types and calculates total cost

## Prerequisites

1. Database must be accessible and configured
2. (Optional) Scraper JSON file path for importing pricing data
   - Example: `../../scripts/scraper/www/instances.json`

## Running the Test

### Option 1: Run Migration First

```bash
# Run migration to add new fields
go run cmd/run_migration/main.go

# Then run the test usecase
go run cmd/platform/main.go
```

### Option 2: Run with Import

```go
import (
    "context"
    "github.com/mo7amedgom3a/arch-visualizer/backend/pkg/usecases/scenario11_ec2_pricing_import"
)

// With scraper JSON file
err := scenario11_ec2_pricing_import.EC2PricingImportTestRunner(
    context.Background(),
    "../../scripts/scraper/www/instances.json",
)

// Without import (just test existing DB rates)
err := scenario11_ec2_pricing_import.EC2PricingImportTestRunner(
    context.Background(),
    "",
)
```

## What It Tests

### 1. Database Migration
- Verifies that `instance_type` and `operating_system` columns are added
- Creates necessary indexes for efficient lookups

### 2. Pricing Import
- Parses scraper JSON format
- Converts EC2 pricing data to `PricingRate` records
- Upserts rates to database with proper deduplication

### 3. Pricing Calculation
Tests pricing for:
- Different instance types: `t3.micro`, `t3.small`, `m5.large`, `c5.xlarge`
- Different regions: `us-east-1`, `us-west-2`
- Different operating systems: `linux`, `mswin` (Windows)

### 4. Rate Source Detection
- Shows whether rates come from database (imported) or fallback (hardcoded)
- Helps verify that import was successful

## Expected Output

```
====================================================================================================
SCENARIO 11: EC2 Pricing Import and Calculation Test
====================================================================================================

[Step 1] Running database migrations...
✓ Database migrations completed successfully

[Step 2] Importing EC2 pricing data from: www/instances.json
✓ EC2 pricing data imported successfully
  Total Instances Processed: 500+
  Total Rates Imported: 5000+
  Regions: 20+
  Operating Systems: 2

[Step 3] Testing pricing calculation with different EC2 instance types...

📊 Pricing Calculation Results:
----------------------------------------------------------------------------------------------------
Instance Type       Region          OS         Hourly Rate     Monthly Cost    Source
----------------------------------------------------------------------------------------------------
t3.micro            us-east-1       linux      $0.010400      $7.49          DB
t3.small            us-east-1       linux      $0.020800      $14.98          DB
m5.large            us-east-1       linux      $0.096000      $69.12          DB
c5.xlarge           us-east-1       linux      $0.170000      $122.40         DB
t3.micro            us-east-1       mswin      $0.020800      $14.98          DB
t3.micro            us-west-2       linux      $0.010400      $7.49           DB

[Step 4] Creating test architecture with multiple EC2 instance types...
✓ Architecture processed successfully
  Project ID: <uuid>

[Step 5] Pricing Breakdown:
====================================================================================================
💰 TOTAL MONTHLY COST: $XX.XX USD
   Provider: aws
   Region: us-east-1
   Period: monthly

📋 Resource-by-Resource Breakdown:
...
```

## Troubleshooting

### Migration Fails
- Ensure database is accessible
- Check that previous migrations have been applied
- Verify database user has ALTER TABLE permissions

### Import Fails
- Verify scraper JSON file path is correct
- Check that JSON file is valid and follows expected format
- Ensure database connection is established

### No Rates Found
- Verify that import completed successfully
- Check that instance types in test match those in imported data
- Ensure region and OS match imported data

### Rates Show "Fallback"
- This means DB lookup failed and hardcoded rates are used
- Verify that import was successful
- Check that instance type, region, and OS match imported data exactly

## Next Steps

After running this test:
1. Verify that imported rates are accurate
2. Test with your own scraper output
3. Extend to other AWS services (RDS, ElastiCache, etc.)
4. Add support for Reserved Instances and Spot pricing
