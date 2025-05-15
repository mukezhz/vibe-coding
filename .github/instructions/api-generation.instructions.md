# API Generation Prompts

Use these structured prompts when asking LLMs like GitHub Copilot to generate REST API components for our Go Clean Architecture project.

## Complete API Feature

```
Generate a complete REST API for a [feature_name] feature with the following requirements:

1. Model definition in domain/models/[feature_name].go with these fields:
   - ID (UUID)
   - [Add specific fields with types and constraints]
   - Created/Updated timestamps

2. All required clean architecture components:
   - Repository with CRUD operations
   - Service with business logic
   - Controller with HTTP handlers
   - DTO structs for requests and responses
   - Routes for API endpoints
   - Module for dependency injection
   - Custom errors (if needed)

3. Implement these API endpoints:
   - POST /api/[feature_name] - Create
   - GET /api/[feature_name]/:id - Get by ID
   - GET /api/[feature_name] - List with pagination
   - PUT /api/[feature_name]/:id - Update
   - DELETE /api/[feature_name]/:id - Delete (optional)

4. Include proper validation, error handling, and follow our project's patterns for:
   - Response formatting
   - Error responses
   - UUID handling with types.BinaryUUID
   - Pagination
   - Validation using Gin binding tags

Follow the architectural patterns shown in domain/todo and domain/organization examples.
```

## Model Definition Only

```
Generate a GORM model definition for [feature_name] with the following fields:

- [List fields with types and descriptions]

Follow these requirements:
- Use proper GORM tags for MySQL compatibility
- Include types.BinaryUUID for UUID fields
- Add BeforeCreate hook for UUID generation if needed
- Follow our naming conventions (UpperCamelCase for struct, snake_case for table)
- Include CreatedAt, UpdatedAt, and DeletedAt fields

Reference the Todo model structure in domain/models/todo.go for:
- Field naming conventions
- GORM tag usage
- Primary key definition
- Time fields handling
```

## Controller Implementation

```
Generate a controller implementation for the [feature_name] feature with the following endpoints:

- [List specific endpoints needed]

Requirements:
- Follow our Controller struct pattern with logger, service, and env dependencies
- Return *Controller as a pointer from NewController constructor
- Implement proper request binding and validation
- Use responses.DetailResponse and responses.ListResponse for consistent responses
- Use responses.HandleError and responses.HandleValidationError for handling errors
- Implement pagination for list endpoints (page and limit query parameters)
- Convert between domain models and DTO formats
- Return appropriate HTTP status codes for different operations

Follow the pattern in domain/todo/controller.go with:
- CreateTodo: Binding request, validating input, creating entity, returning response
- GetTodoByID: Parsing ID parameter, fetching entity, returning detailed response
- UpdateTodo: Parsing ID, fetching entity, binding request, updating specific fields
- FetchTodoWithPagination: Handling pagination parameters, converting to list response
```

## Repository Layer

```
Create a repository implementation for [feature_name] with these operations:

- [List operations needed, e.g., Create, GetByID, Update, Delete, List]

Include:
- Return Repository as a value (not a pointer) from NewRepository constructor
- Proper error handling using GORM error patterns
- Logging for all operations with context information
- Transaction support where needed
- Efficient pagination implementation for List operation:
  - Count total records
  - Apply offset and limit
  - Return items and total count
- Any specific queries/filters needed:
  [List specific filters or query needs]

Follow the pattern in domain/todo/repository.go with:
- Repository struct embedding infrastructure.Database
- Consistent method signatures
- Clean implementation of database operations
```

## Service Layer

```
Generate a service layer for [feature_name] with business logic for:

- [List operations needed: Create, GetByID, Update, List, etc.]

Requirements:
- Follow our Service struct pattern with logger and repository dependencies
- Return *Service as a pointer from NewService constructor
- Implement proper error mapping (e.g., convert 'record not found' to domain-specific errors)
- Include validation and business rules:
  [List specific rules or validations]
- Add proper logging and error handling
- Implement any complex business logic:
  [Describe any complex logic needed]

Follow the pattern in domain/todo/service.go with:
- Clean method signatures (Create, GetByID, Update, List)
- Domain-specific error returns
- Simple pass-through to repository for CRUD operations where appropriate
```

## Complete Route Definition

```
Create a route definition for the [feature_name] feature with these endpoints:

- [List HTTP methods and paths]

Follow the project pattern with:
- Route struct with logger, handler, and controller dependencies
- NewRoute constructor returning a pointer
- RegisterRoute function to set up all endpoints
- Proper grouping under /api/[feature_name]

Follow the pattern in domain/todo/route.go with:
- Group setup with handler.Group
- Clear HTTP method definitions (POST, GET, PUT, DELETE)
- Clean mapping to controller methods
```

## Error Definitions

```
Create custom error definitions for the [feature_name] feature:

Error cases to cover:
- [List error scenarios like "not found", "already exists", etc.]

Follow our error pattern with:
- Error constants (e.g., ERR_FEATURE_NOT_FOUND)
- Constructor functions (e.g., NewFeatureNotFoundError)
- Appropriate HTTP status codes

Follow the pattern in domain/todo/errorz.go with:
- Consistent error codes in UPPERCASE_WITH_UNDERSCORES format
- Constructor functions that return properly formatted errors
- Integration with the pkg/errorz package for consistent error handling
```

## DTO Definitions

```
Generate DTOs (Data Transfer Objects) for the [feature_name] feature:

Required DTOs:
- CreateRequest with validation tags
- UpdateRequest with validation tags and pointer fields
- Response format
- ListItem format for list responses
- ListResponse with pagination

Include proper validation using binding tags:
- [List specific validation needs]
```

## Module Setup

```
Generate the module setup for [feature_name] with fx dependency injection:

- Provide all dependencies in correct order
- Invoke route registration
- Follow the pattern in other feature modules
```
