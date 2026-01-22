# Project Backlog

This document tracks completed features, work in progress, and planned features for the Architecture Visualizer project.

## How to Use This Backlog

- **Completed**: Features that are fully implemented and tested
- **In Progress**: Features currently being worked on
- **Planned**: Features planned for future implementation

When updating this backlog:
- Move items from "Planned" to "In Progress" when work begins
- Move items from "In Progress" to "Completed" when work is finished
- Add new planned items as they are identified
- Include brief descriptions and references to relevant documentation

---

## Completed

### Core Architecture
- ✅ Modular monolith architecture with clear separation of concerns
- ✅ Domain-first design with cloud-agnostic core
- ✅ Strong module boundaries and provider-specific implementations
- ✅ Rule-driven validation system (no hardcoded logic)
- ✅ Pluggable IaC engines architecture

### AWS Provider Implementation
- ✅ AWS networking resources (VPC, Subnet, Internet Gateway, Route Tables, Security Groups, NAT Gateway)
- ✅ AWS compute resources (EC2, Lambda, Application Load Balancer, Target Groups, Listeners, Auto Scaling Groups, Launch Templates)
- ✅ AWS storage resources (S3 Buckets, EBS Volumes)
- ✅ AWS IAM resources (IAM Roles, IAM Instance Profiles)
- ✅ AWS adapter pattern implementation for networking, compute, and storage
- ✅ AWS mapper pattern for domain ↔ cloud model conversion
- ✅ AWS rules implementation using domain rules interfaces
- ✅ AWS SDK integration layer

### Infrastructure as Code (IaC)
- ✅ Terraform engine implementation
- ✅ Pulumi engine implementation
- ✅ Code generation orchestration pipeline
- ✅ Resource dependency resolution
- ✅ Multi-file IaC output generation

### Domain Layer
- ✅ Cloud-agnostic domain models for resources
- ✅ Resource relationships and dependencies
- ✅ Architecture aggregates and graph representation
- ✅ Domain rules system with interfaces
- ✅ Constraint evaluator and validation logic

### Diagram & Parsing
- ✅ Canvas JSON parser
- ✅ Graph representation builder
- ✅ Structural validation (cyclic containment, connector types, parent/child relationships)
- ✅ Domain-ready architecture preparation

### Rules Engine
- ✅ Database-driven validation rules
- ✅ Domain rules interfaces (cloud-agnostic)
- ✅ AWS-specific rule implementations
- ✅ Constraint types: requires_parent, allowed_parent, requires_region, max_children, min_children

### Pricing Service
- ✅ AWS networking resources pricing
- ✅ Flexible pricing models (per hour, per GB, per request)
- ✅ Pricing calculator and service interfaces
- ✅ Cost estimation for complete architectures

### Use Cases & Examples
- ✅ Scenario 1: Basic Web Application (3-tier architecture)
- ✅ Scenario 2: High Availability Architecture (multi-AZ with load balancing)
- ✅ Scenario 3: Scalable API Architecture (auto-scaling with IAM)
- ✅ Scenario 4: Lambda + S3 Integration (serverless architecture)
- ✅ Mock helpers for all resource types
- ✅ Region selection and validation utilities

### API Layer
- ✅ HTTP endpoints and routing
- ✅ Request/response DTOs
- ✅ Middleware support
- ✅ Authentication & authorization structure

### Data Access
- ✅ PostgreSQL persistence layer
- ✅ Repository pattern implementation
- ✅ Database schema for projects, resources, constraints, relationships

### Documentation
- ✅ Backend architecture documentation
- ✅ Workflow documentation (frontend ↔ backend)
- ✅ AWS adapters documentation
- ✅ Use cases documentation
- ✅ Domain rules documentation
- ✅ Pricing service documentation

---

## In Progress

_No items currently in progress. Update this section as work begins on new features._

---

## Planned

### Additional Cloud Providers
- 🔲 GCP provider implementation (models, mappers, adapters, services)
- 🔲 Azure provider implementation (models, mappers, adapters, services)
- 🔲 Multi-provider architecture support

### Additional IaC Engines
- 🔲 AWS CDK engine implementation
- 🔲 CloudFormation template generation
- 🔲 Ansible playbook generation

### Compute Resources
- 🔲 ECS (Elastic Container Service) support
- 🔲 EKS (Elastic Kubernetes Service) support
- 🔲 EC2 additional features:
  - 🔲 EBS volume attachments (beyond root volume)
  - 🔲 Key Pairs validation and management
  - 🔲 Placement Groups support
  - 🔲 Spot Instances support

### Serverless Resources
- 🔲 API Gateway integration
- 🔲 Lambda function invocations pricing
- 🔲 EventBridge (CloudWatch Events) support
- 🔲 Step Functions support

### Database Resources
- 🔲 RDS (Relational Database Service) support
- 🔲 DynamoDB support
- 🔲 ElastiCache support
- 🔲 Database pricing calculations

### Networking Resources
- 🔲 CloudFront CDN integration
- 🔲 VPN Gateway support
- 🔲 Direct Connect support
- 🔲 Transit Gateway support

### Storage Resources
- 🔲 S3 advanced features (lifecycle policies, replication)
- 🔲 EFS (Elastic File System) support
- 🔲 Glacier support

### Monitoring & Observability
- 🔲 CloudWatch integration
- 🔲 CloudTrail integration
- 🔲 Monitoring dashboard generation
- 🔲 Alert configuration

### Cost Management
- 🔲 Compute resources pricing (EC2, Lambda, ECS)
- 🔲 Serverless resources pricing (API Gateway, Lambda invocations)
- 🔲 Storage resources pricing (S3, EBS detailed pricing)
- 🔲 Database resources pricing (RDS, DynamoDB)
- 🔲 Dynamic pricing via AWS Pricing API integration
- 🔲 Cost tracking for created resources
- 🔲 Cost alerts and thresholds
- 🔲 Cost optimization suggestions
- 🔲 Multi-provider pricing support (GCP, Azure)

### Architecture Features
- 🔲 Multi-region architecture support
- 🔲 Disaster recovery scenarios
- 🔲 Architecture validation rules expansion
- 🔲 Visual diagram generation from architecture
- 🔲 Architecture templates library
- 🔲 Best practices recommendations

### Frontend Integration
- 🔲 Real-time validation feedback
- 🔲 Code preview in UI
- 🔲 One-click deployment integration
- 🔲 Git repository push functionality
- 🔲 In-editor visualization

### Advanced Features
- 🔲 Architecture versioning
- 🔲 Architecture comparison and diff
- 🔲 Export to multiple formats (JSON, YAML, HCL)
- 🔲 Import existing Terraform/Pulumi code
- 🔲 Architecture templates and presets
- 🔲 Collaborative editing support

### Testing & Quality
- 🔲 Comprehensive integration tests
- 🔲 End-to-end testing framework
- 🔲 Performance testing and optimization
- 🔲 Load testing for API endpoints

### Documentation
- 🔲 API documentation (OpenAPI/Swagger)
- 🔲 Developer guide
- 🔲 Deployment guide
- 🔲 Contributing guidelines

---

## Notes

- Items are organized by category for easier tracking
- Priority and timeline information can be added to individual items as needed
- Reference specific documentation files when relevant (e.g., `backend/README.md`, `backend/workflow.md`)
- Update this backlog regularly to reflect current project status
