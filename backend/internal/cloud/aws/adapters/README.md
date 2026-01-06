# AWS Adapters Layer

This directory contains adapter implementations that bridge the domain layer and AWS-specific services.

## Purpose

The adapter layer implements the **Adapter Pattern** to:
- Decouple domain layer from cloud-specific implementations
- Enable easy provider swapping
- Provide a clean abstraction for domain services
- Handle conversion between domain and cloud models

## Structure

```
adapters/
└── networking/
    ├── adapter.go          # Main adapter implementation
    ├── adapter_test.go     # Adapter tests
    ├── factory.go          # Factory pattern for adapter creation
    └── README.md           # Detailed documentation
```

## Architecture Flow

```
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer                            │
│  (Cloud-Agnostic Business Logic)                            │
│                                                             │
│  NetworkingService Interface                               │
└──────────────────────┬────────────────────────────────────┘
                        │
                        │ implements
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│                  Adapter Layer                              │
│  (This Package)                                             │
│                                                             │
│  AWSNetworkingAdapter                                       │
│  - Implements NetworkingService                             │
│  - Uses Mappers for conversion                              │
│  - Handles validation at both layers                        │
└──────────────────────┬────────────────────────────────────┘
                        │
                        │ uses
                        ▼
┌─────────────────────────────────────────────────────────────┐
│              AWS Service Layer                              │
│  (Cloud-Specific Implementation)                            │
│                                                             │
│  AWSNetworkingService Interface                            │
└──────────────────────┬────────────────────────────────────┘
                        │
                        │ uses
                        ▼
┌─────────────────────────────────────────────────────────────┐
│              AWS Models & API                               │
│  (AWS SDK, Terraform, etc.)                                 │
└─────────────────────────────────────────────────────────────┘
```

## Key Benefits

1. **Separation of Concerns**: Domain layer never imports AWS code
2. **Testability**: Easy to mock AWS services
3. **Extensibility**: Add new providers without changing domain code
4. **Type Safety**: Compile-time interface compliance
5. **Error Context**: Clear error messages with full context

## Usage Example

```go
// 1. Create AWS service implementation
awsService := &awsNetworkingServiceImpl{...}

// 2. Create adapter
adapter := networking.NewAWSNetworkingAdapter(awsService)

// 3. Use as domain service (cloud-agnostic)
var service domainnetworking.NetworkingService = adapter

// 4. Domain layer uses service without knowing about AWS
vpc, err := service.CreateVPC(ctx, domainVPC)
```

## Adding New Adapters

To add adapters for other resource types:

1. Create new directory: `adapters/{resource_type}/`
2. Implement domain service interface
3. Use mappers for conversion
4. Add factory pattern
5. Write tests

Example: `adapters/compute/` for EC2, Lambda, etc.

## Future Providers

The adapter pattern makes it easy to add new cloud providers:

1. Create `internal/cloud/{provider}/adapters/networking/`
2. Implement `domainnetworking.NetworkingService`
3. Use provider-specific service and mappers
4. Add to factory pattern

The domain layer remains completely unchanged!
