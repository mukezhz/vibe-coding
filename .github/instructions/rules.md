# copilot rules - Go Clean Architecture

## Project Patterns

### Domain Organization
- Each domain has consistent files: controller.go, service.go, repository.go, route.go, module.go, dto.go, errorz.go
- Models are centralized in domain/models with GORM annotations
- Each domain has its own module with dependency injection setup
- Route registration is invoked through fx.Invoke(RegisterRoute)

### Naming Conventions
- Use PascalCase for exported functions and types
- Controller methods are named for their actions (CreateTodo, GetTodoByID, etc.)
- Repository methods mirror service methods (Create, GetByID, etc.)
- Error types use descriptive names with "Error" suffix (e.g., TodoNotFoundError)

### Error Handling
- Custom errors defined in domain/*/errorz.go
- Standard error handling with responses.HandleError and responses.HandleValidationError
- Common errors defined in pkg/errorz
- HTTP status codes mapped to error types

### Response Formatting
- Standard response types in pkg/responses
- DetailResponse for single item responses
- ListResponse for paginated list responses
- Error responses with consistent structure

### Data Types
- Binary UUID for primary keys with custom type
- Consistent time handling with time.Time
- DTO structs for request/response mapping

### Middleware Usage
- Authentication middleware for protected routes
- Rate limiting with ulule/limiter
- Global middleware registered in pkg/middlewares

### AWS Integration
- S3 service for file storage
- Cognito for authentication
- SES for email sending
- Environment variables for AWS configuration

### Database Operations
- GORM for database interactions
- Repository pattern for data access
- MySQL dialect configured
- Atlas for migrations

### Command Structure
- Cobra for CLI commands
- app:serve as main command
- Command registration in console/console.go

## Critical Implementation Paths

### Request Flow
1. Request → Gin Router
2. Router → Middleware
3. Middleware → Route Handler
4. Route Handler → Controller
5. Controller → Service
6. Service → Repository
7. Repository → Database
8. Response back through layers

### New Domain Addition
1. Create domain folder
2. Implement model in domain/models
3. Create controller, service, repository, route
4. Create module with dependency injection
5. Register module in domain/module.go
6. Create migration files if needed

### Database Changes
1. Update model in domain/models
2. Run make migrate-diff to generate migration
3. Run make migrate-apply to apply changes

### Authentication Flow
1. Token extraction from Authorization header
2. Verification with Cognito service
3. Claims extraction and context enrichment
4. Role-based access control in routes

## Known Challenges
- AWS service integration requires proper configuration
- Error handling might need enhancement for specific cases
- Database migrations need careful management
- Authentication with Cognito needs proper setup

## Preferred Tools
- Bruno for API testing
- Docker for development environment
- Make for common commands
- Atlas for database migrations

## Testing Approach
- Integration tests appear to be the main focus
- Unit tests for specific components
- Mock dependencies when needed
