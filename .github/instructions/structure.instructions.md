# GitHub Copilot Instructions for Go Clean Architecture Project

## Project Overview
This project follows the Clean Architecture principles using Go and Gin Web Framework. The codebase is structured to separate concerns and maintain a clear hierarchy of dependencies.
The codebase follow pragmatic approach to clean architecture, focusing on simplicity and ease of understanding. The goal is to create a maintainable and scalable codebase that adheres to best practices in Go development.
The project uses less interface and more concrete types, making it easier to understand and work with. The focus is on creating a clean and organized structure that allows for easy navigation and understanding of the codebase.

## Architecture & Directory Structure

- `/bootstrap`: Modules required to start the application
- `/console`: Server commands
- `/domain`: Contains models, constants, and domain-specific code
  - `/domain/constants`: Global application constants
  - `/domain/models`: ORM models
  - `/domain/<domain-name>`: Domain-specific route, controller, sevice and repository
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
- `/tests`: Application tests

## Coding Conventions

### Dependency Injection
- Use uber-go/fx for dependency injection
- Register all modules in bootstrap/app.go
- Follow the pattern of providing components and invoking functions

### Domain Structure
For each domain (e.g., user, product, order), follow this structure:
- controller.go: Handles HTTP requests and responses
- repository.go: Database access layer
- route.go: Defines API endpoints
- service.go: Business logic
- dto.go: Data Transfer Object
- errorz.go: Custom Error type

### Error Handling
- Use custom error types from pkg/errorz
- Return structured errors with proper HTTP status codes
- Follow the error handling pattern used in the codebase

### Database Operations
- Use GORM for database operations
- Follow the repository pattern
- Use models defined in domain/models

## Common Patterns

### Creating a New Domain
1. Create a new directory under `/domain` (e.g., `/domain/product`)
2. Create controller, repository, routes, and service files
3. Register the domain module in bootstrap/app.go

### Database Migration

Below are the supported `make` commands for managing database migrations:

| Make Command          | Description                                                                 |
| --------------------- | --------------------------------------------------------------------------- |
| `make migrate-status` | Show the current migration status                                           |
| `make migrate-diff`   | Generate a new migration by comparing models to the current DB (`gorm` env) |
| `make migrate-apply`  | Apply all pending migrations                                                |
| `make migrate-down`   | Roll back the most recent migration (`gorm` env)                            |
| `make migrate-hash`   | Hash migration files for integrity checking      

