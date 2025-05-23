# System Patterns - Go Clean Architecture

## System Architecture
This application follows a clean architecture approach with clear separation of concerns:

```mermaid
flowchart TD
    subgraph Presentation
        RT[Route Layer]
        CT[Controller Layer]
    end
    
    subgraph Business
        SV[Service Layer]
    end
    
    subgraph Data
        RP[Repository Layer]
        MD[Models]
    end
    
    subgraph Infrastructure
        DB[Database]
        RO[Router]
        AWS[AWS Services]
        MW[Middlewares]
    end
    
    RT --> CT
    CT --> SV
    SV --> RP
    RP --> MD
    MD --> DB
    RT --> RO
    SV --> AWS
    RT --> MW
```

## Domain-Driven Structure
Each domain (feature area) is organized with the same internal structure:

```mermaid
flowchart LR
    subgraph Domain[Domain Package]
        C[controller.go]
        S[service.go]
        R[repository.go]
        RT[route.go]
        M[module.go]
        DTO[dto.go]
        E[errorz.go]
    end
    
    M --> C & S & R & RT
    C --> DTO
    C --> E
    S --> R
    RT --> C
```

## Key Technical Decisions

### 1. Dependency Injection with uber-go/fx
- Used for managing dependencies throughout the application
- Modules are registered in the bootstrap/modules.go file
- Each domain has its own module.go to register dependencies

### 2. Command Structure with Cobra
- CLI commands structured with Cobra
- Main app:serve command for running the API server
- Shared application container for commands

### 3. Database Access with GORM
- GORM for database operations
- Repository pattern for data access
- Models with GORM annotations in domain/models
- Atlas for database migrations

### 4. HTTP Handling with Gin
- Gin Web Framework for HTTP routing
- Standardized response structures
- Middleware support for authentication, rate limiting, etc.

### 5. Environment Configuration
- Environment variables managed with Viper
- Central framework.Env structure
- .env file support with godotenv

### 6. Error Handling
- Custom error types with proper HTTP status codes
- Standardized error responses
- Error handling middleware with Sentry integration

## Component Relationships

### Application Bootstrap
```mermaid
flowchart TD
    Main[main.go] --> Bootstrap[bootstrap/app.go]
    Bootstrap --> Console[console/console.go]
    Bootstrap --> Modules[bootstrap/modules.go]
    Modules --> PKG[pkg/module.go]
    Modules --> Domain[domain/module.go]
    Modules --> Seeds[seeds/module.go]
    Domain --> UserModule[domain/user/module.go]
    Domain --> TodoModule[domain/todo/module.go]
    Domain --> OrgModule[domain/organization/module.go]
    Domain --> BookingModule[domain/booking/module.go]
    PKG --> Framework[pkg/framework]
    PKG --> Infrastructure[pkg/infrastructure]
    PKG --> Middlewares[pkg/middlewares]
    PKG --> Services[pkg/services]
```

### Request Flow
```mermaid
flowchart LR
    Request[Request] --> Router[Gin Router]
    Router --> Middleware[Middlewares]
    Middleware --> Route[Domain Route]
    Route --> Controller[Domain Controller]
    Controller --> Service[Domain Service]
    Service --> Repository[Domain Repository]
    Repository --> Database[GORM/Database]
    Controller --> Response[Response]
```

### Dependency Injection
```mermaid
flowchart TD
    Module[fx Module] -->|Provides| Components[Components]
    Module -->|Invokes| Functions[Functions]
    Components --> Controller[Controllers]
    Components --> Service[Services]
    Components --> Repository[Repositories]
    Components --> Route[Routes]
    Functions --> RegisterRoutes[Register Routes]
    Functions --> Migrate[Migrate]
```

The architecture follows a consistent pattern across all domains, making it scalable and maintainable. The dependency injection system ensures loose coupling between components.
