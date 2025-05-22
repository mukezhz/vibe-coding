---
applyTo: "**/*.go"
---
# Project Structure and Architecture Guide

This guide provides a comprehensive overview of the project structure, architecture, and best practices for our Go Clean Architecture project.

## Project Overview

This project follows Clean Architecture principles using Go and the Gin Web Framework. The codebase is structured to separate concerns and maintain a clear hierarchy of dependencies. We take a pragmatic approach to clean architecture, focusing on simplicity and ease of understanding.

## Directory Structure

- `/bootstrap`: Modules required to start the application
- `/console`: Server commands
- `/domain`: Contains models, constants, and domain-specific code
  - `/domain/constants`: Global application constants
  - `/domain/models`: ORM models
  - `/domain/<domain-name>`: Domain-specific components (controller, service, repository, etc.)
- `/pkg`: Shared packages
  - `/pkg/errorz`: Custom error types and handlers
  - `/pkg/framework`: Core framework components
  - `/pkg/infrastructure`: Third-party service connections
  - `/pkg/middlewares`: HTTP request middlewares
  - `/pkg/responses`: Standardized HTTP response structures
  - `/pkg/services`: Shared application services
  - `/pkg/types`: Custom data types
  - `/pkg/utils`: Global utility functions
- `/migrations`: Database migration files managed by Atlas
- `/testutil`: Testing utilities

## Clean Architecture Layers

Each domain follows a consistent layer structure:

1. **Models Layer** (`domain/models/`): Database entities with GORM annotations
2. **Repository Layer** (`domain/<feature>/repository.go`): Database access logic
3. **Service Layer** (`domain/<feature>/service.go`): Business logic
4. **Controller Layer** (`domain/<feature>/controller.go`): Request handling and response formation
5. **Route Layer** (`domain/<feature>/route.go`): Endpoint definitions
6. **DTO Layer** (`domain/<feature>/dto.go`): Data Transfer Objects for requests/responses
7. **Error Layer** (`domain/<feature>/errorz.go`): Feature-specific custom errors
8. **Module Layer** (`domain/<feature>/module.go`): Dependency injection via fx

## Dependency Injection

We use the `uber-go/fx` library for dependency injection:

```go
// domain/feature/module.go
package feature

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewRepository),
	fx.Provide(NewService),
	fx.Provide(NewController),
	fx.Provide(NewRoute),
	fx.Invoke(func(route *Route) {
		route.RegisterRoute()
	}),
)
```

Dependency registration in bootstrap:

```go
// bootstrap/modules.go
package bootstrap

import (
	"clean-architecture/domain/todo"
	"clean-architecture/domain/organization"
	// Other imports

	"go.uber.org/fx"
)

var DomainModules = fx.Options(
	todo.Module,
	organization.Module,
	// Other domain modules
)
```

## Database Operations

The project uses GORM for database operations, following the repository pattern:

```go
// domain/todo/repository.go
package todo

import (
	"clean-architecture/domain/models"
	"clean-architecture/pkg/infrastructure"
	"clean-architecture/pkg/types"
)

type Repository struct {
	infrastructure.Database
}

func NewRepository(db infrastructure.Database) Repository {
	return Repository{db}
}

func (r Repository) Create(todo *models.Todo) error {
	return r.DB().Create(todo).Error
}

// Other repository methods...
```

## Database Migration

The project uses Atlas for database migrations:

| Make Command          | Description                                 |
| --------------------- | ------------------------------------------- |
| `make migrate-status` | Show current migration status               |
| `make migrate-diff`   | Generate new migration                      |
| `make migrate-apply`  | Apply all pending migrations                |
| `make migrate-down`   | Roll back the most recent migration         |
| `make migrate-hash`   | Hash migration files for integrity checking |

## Go Coding Guidelines

### Code Organization
- Follow the wesionaryTEAM/go_clean_architecture Go project layout
- Use packages to organize code by functionality, not by type
- Keep package names simple, short, and meaningful
- One package per directory

### Naming Conventions
- Use `camelCase` for private variables, functions, and methods
- Use `PascalCase` for exported (public) variables, functions, methods, and types
- Use short but descriptive names for variables
- Prefer clarity over brevity for function and method names
- Use acronyms consistently (e.g., `HTTP`, `URL`, `ID`)

### Error Handling
- Check errors immediately after function calls
- Don't use `panic` or `recover` in production code
- Return errors rather than using panic
- Use custom error types for specific error conditions
- Use `errors.Is()` and `errors.As()` functions for error checking
- Wrap errors with context using `fmt.Errorf("doing X: %w", err)`
