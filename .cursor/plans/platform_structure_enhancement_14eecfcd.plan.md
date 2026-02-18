---
name: ""
overview: ""
todos: []
---

# Platform Structure Enhancement Plan

## Current State Analysis

The platform package currently has:

- **30+ model files** in a flat structure (`models/`)
- **30+ repository files** in a flat structure (`repository/`)
- **19 service files** split between `server/services/` and `services/`
- **Test files** mixed with implementation files
- **Inconsistent naming** (e.g., `project_service.go` vs `project_service_architecture.go`)
- **Repository adapters** located in services package

## Identified Business Domains

1. **Project Domain**: Project lifecycle, versions, variables, outputs, UI state
2. **Resource Domain**: Resources, types, categories, kinds, constraints, relationships (containment, dependencies)
3. **Template/Marketplace Domain**: Templates, categories, reviews, technologies, compliance standards, IAC formats
4. **Pricing Domain**: Pricing rates, components, project/resource/service pricing
5. **User Domain**: User management
6. **Infrastructure Domain**: IAC targets and lookup tables

## Proposed Structure

```
backend/internal/platform/
├── models/
│   ├── project/          # Project domain models
│   │   ├── project.go
│   │   ├── project_version.go
│   │   ├── project_variable.go
│   │   ├── project_output.go
│   │   └── ui_state.go
│   ├── resource/         # Resource domain models
│   │   ├── resource.go
│   │   ├── resource_type.go
│   │   ├── resource_category.go
│   │   ├── resource_kind.go
│   │   ├── resource_constraint.go
│   │   ├── resource_containment.go
│   │   ├── resource_dependency.go
│   │   ├── dependency_type.go
│   │   └── hidden_dependency.go
│   ├── template/          # Template/Marketplace domain models
│   │   ├── template.go
│   │   ├── template_component.go
│   │   ├── template_feature.go
│   │   ├── template_use_case.go
│   │   ├── template_compliance.go
│   │   ├── template_technology.go
│   │   ├── template_iac_format.go
│   │   ├── category.go
│   │   ├── review.go
│   │   ├── technology.go
│   │   ├── compliance_standard.go
│   │   └── iac_format.go
│   ├── pricing/           # Pricing domain models
│   │   ├── pricing_rate.go
│   │   ├── pricing_component.go
│   │   ├── project_pricing.go
│   │   ├── resource_pricing.go
│   │   ├── service_pricing.go
│   │   └── service_type_pricing.go
│   ├── user/              # User domain models
│   │   └── user.go
│   └── infrastructure/    # Infrastructure lookup models
│       └── iac_target.go
│
├── repository/
│   ├── base.go            # Base repository (shared)
│   ├── project/           # Project domain repositories
│   │   ├── project_repository.go
│   │   ├── project_version_repository.go
│   │   ├── project_variable_repository.go
│   │   └── project_output_repository.go
│   ├── resource/          # Resource domain repositories
│   │   ├── resource_repository.go
│   │   ├── resource_type_repository.go
│   │   ├── resource_category_repository.go
│   │   ├── resource_kind_repository.go
│   │   ├── resource_constraint_repository.go
│   │   ├── resource_containment_repository.go
│   │   ├── resource_dependency_repository.go
│   │   ├── dependency_type_repository.go
│   │   └── hidden_dependency_repository.go
│   ├── template/          # Template domain repositories
│   │   ├── template_repository.go
│   │   ├── template_compliance_repository.go
│   │   ├── template_technology_repository.go
│   │   ├── template_iac_format_repository.go
│   │   ├── category_repository.go
│   │   ├── review_repository.go
│   │   ├── technology_repository.go
│   │   ├── compliance_standard_repository.go
│   │   └── iac_format_repository.go
│   ├── pricing/           # Pricing domain repositories
│   │   ├── pricing_repository.go
│   │   └── pricing_rate_repository.go
│   ├── user/              # User domain repositories
│   │   └── user_repository.go
│   ├── infrastructure/    # Infrastructure repositories
│   │   └── iac_target_repository.go
│   └── tests/             # Repository tests (moved from repository/tests/)
│       ├── base_repository_test.go
│       ├── project_repository_test.go
│       ├── resource_repository_test.go
│       ├── pricing_repository_test.go
│       ├── user_repository_test.go
│       ├── marketplace_repositories_test.go
│       └── lookup_repositories_test.go
│
├── server/
│   ├── interfaces/        # Service interfaces (unchanged)
│   ├── orchestrator/      # Orchestration logic (unchanged)
│   ├── services/          # Consolidated service implementations
│   │   ├── project/       # Project domain services
│   │   │   ├── project_service.go
│   │   │   └── project_service_architecture.go
│   │   ├── resource/      # Resource domain services
│   │   │   ├── resource_metadata_service.go
│   │   │   └── constraint_service.go
│   │   ├── template/      # Template domain services (if any)
│   │   ├── pricing/       # Pricing domain services
│   │   │   ├── pricing_service.go
│   │   │   └── pricing_importer/  # Move from platform/services/
│   │   │       ├── importer.go
│   │   │       ├── ec2_parser.go
│   │   │       └── models.go
│   │   ├── user/          # User domain services
│   │   │   └── user_service.go
│   │   ├── core/          # Core orchestration services
│   │   │   ├── diagram_service.go
│   │   │   ├── architecture_service.go
│   │   │   ├── codegen_service.go
│   │   │   ├── optimization_service.go
│   │   │   └── static_data_service.go
│   │   ├── adapters/      # Repository adapters (moved from services/)
│   │   │   └── repository_adapters.go
│   │   └── rules/         # Rule service adapters
│   │       └── rule_service_adapter.go
│   └── tests/             # Service tests (separated)
│       ├── project/
│       │   ├── project_service_test.go
│       │   └── project_service_architecture_test.go
│       ├── resource/
│       │   └── constraint_service_test.go
│       ├── pricing/
│       │   └── pricing_service_test.go
│       └── core/
│           ├── diagram_service_test.go
│           ├── codegen_service_test.go
│           └── architecture_service_test.go
│
├── config/                # Configuration (unchanged)
├── database/              # Database connection (unchanged)
├── logger/                # Logging (unchanged)
├── auth/                  # Authentication (unchanged)
├── errors/                # Error handling (unchanged)
└── utils/                 # Utilities (unchanged)
```

## Implementation Steps

### Phase 1: Models Reorganization

1. Create domain subdirectories under `models/`
2. Move model files to appropriate domain directories
3. Update package declarations from `package models` to `package project`, `package resource`, etc.
4. Update all imports across the codebase to use new package paths
5. Update GORM table names if needed (they should remain the same)

### Phase 2: Repository Reorganization

1. Create domain subdirectories under `repository/`
2. Move repository files to appropriate domain directories
3. Update package declarations
4. Update imports in repositories (especially model imports)
5. Move `repository/tests/` to `repository/tests/` (update package names)
6. Update all service imports to use new repository paths

### Phase 3: Services Consolidation & Reorganization

1. Move `platform/services/pricing_importer/` to `server/services/pricing/pricing_importer/`
2. Create domain subdirectories under `server/services/`
3. Move service files to appropriate domain directories
4. Move `repository_adapters.go` to `server/services/adapters/`
5. Move `rule_service_adapter.go` to `server/services/rules/`
6. Update package declarations
7. Update all imports (especially repository and model imports)
8. Update server initialization code

### Phase 4: Test Separation

1. Create `server/tests/` directory structure matching services
2. Move `*_test.go` files from `server/services/` to `server/tests/`
3. Update test package names and imports
4. Ensure tests can still access internal packages

### Phase 5: Naming Consistency

1. Rename `project_service_architecture.go` to a more descriptive name or merge with `project_service.go`
2. Review all file names for consistency
3. Ensure naming follows Go conventions (snake_case for files, PascalCase for types)

### Phase 6: Documentation Updates

1. Update `models/README.md` with new structure
2. Update `repository/README.md` with new structure
3. Update `server/README.md` with new structure
4. Update `agent.md` files to reflect new organization
5. Update any architecture documentation

## Key Considerations

### Package Naming

- Models: `package project`, `package resource`, `package template`, `package pricing`, `package user`, `package infrastructure`
- Repositories: `package project`, `package resource`, etc. (same as models)
- Services: `package project`, `package resource`, `package core`, `package adapters`, `package rules`

### Import Path Updates

All imports will need to be updated:

- `models.Project` → `project.Project` or `projectmodels "github.com/.../models/project"`
- `repository.ProjectRepository` → `projectrepo "github.com/.../repository/project"`
- Services will need updated imports for both models and repositories

### Backward Compatibility

- This is a breaking change requiring updates across the entire codebase
- Consider creating a migration script or tool to help with import updates
- All tests will need to be updated

### Testing Strategy

- Run tests after each phase to catch issues early
- Update integration tests that may depend on package structure
- Ensure test utilities are accessible from new test locations

## Benefits

1. **Better Discoverability**: Related code is grouped together
2. **Clearer Boundaries**: Domain boundaries are explicit
3. **Easier Maintenance**: Changes to a domain are localized
4. **Scalability**: Easy to add new domains or extend existing ones
5. **Reduced Cognitive Load**: Developers can focus on one domain at a time
6. **Better Testing**: Tests are organized alongside code but separated for clarity

## Migration Risks & Mitigation

1. **Risk**: Breaking imports across large codebase

   - **Mitigation**: Use IDE refactoring tools, create search/replace scripts, test incrementally

2. **Risk**: Circular dependencies between domains

   - **Mitigation**: Review dependencies before moving, use interfaces where needed

3. **Risk**: Test failures due to package visibility

   - **Mitigation**: Keep tests in separate packages but ensure they can access internal types

4. **Risk**: GORM relationships may break

   - **Mitigation**: GORM uses table names, not package names, so relationships sho