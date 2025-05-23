# Active Context - Go Clean Architecture

## Current Work Focus
The current focus is on setting up and understanding the clean architecture Go template. This includes:

1. Comprehending the overall structure and architecture of the application
2. Understanding the domain-driven design approach
3. Learning the flow of requests through the system
4. Exploring the dependency injection patterns with Uber FX
5. Analyzing the error handling mechanisms
6. Understanding the database interaction with GORM and Atlas migrations

## Recent Changes
Initial exploration of the codebase shows:

1. Four functional domains implemented: `user`, `todo`, `organization`, and `booking`
2. Basic CRUD operations in each domain
3. Database setup with MySQL
4. AWS service integration for S3, Cognito, and SES
5. Middleware implementation for authentication, rate limiting, and error handling

## Next Steps
1. Complete the exploration of the codebase and its architecture
2. Document key findings in the memory bank
3. Set up a local development environment
4. Run the application and test API endpoints
5. Potentially add new domains or enhance existing ones
6. Improve testing coverage

## Active Decisions and Considerations

### Architecture Decisions
- The architecture follows a clean, pragmatic approach that prioritizes maintainability over strict adherence to theoretical patterns
- Domain-centric organization with standard components per domain (controller, service, repository, routes)
- Dependency injection with Uber FX for loose coupling
- Standardized error handling and responses

### Technical Considerations
- MySQL is used as the primary database with GORM as the ORM
- Atlas handles database migrations
- Gin Web Framework for HTTP routing and middleware
- AWS SDK integration for cloud services
- JWT-based authentication with Cognito
- Standardized response formats for consistency

### Development Workflow
- Local development can be done directly with Go or through Docker
- Database migrations are managed through Makefile commands and Atlas
- API documentation is provided with Bruno client

### Current Questions
1. What are the specific authentication flows implemented?
2. How are database transactions handled across repositories?
3. What is the deployment strategy for production?
4. How are tests structured and what is the current coverage?
5. Are there any performance bottlenecks to be aware of?
6. What monitoring and observability tools are integrated?

These active considerations will guide further exploration and understanding of the codebase.
