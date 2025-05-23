# Technical Context - Go Clean Architecture

## Technologies Used

### Core Technologies
- **Go** (1.22.4) - Main programming language
- **Gin** - Web framework for HTTP handling
- **GORM** - ORM for database operations
- **Uber FX** - Dependency injection framework
- **MySQL** - Primary database
- **Atlas** - Database migration tool
- **Cobra** - CLI command structure
- **Viper** - Configuration management
- **Zap** - Structured logging
- **AWS SDK** - Integration with AWS services
- **JWT** - Authentication with JWX library
- **Sentry** - Error tracking and monitoring

### Development Tools
- **Docker** - Containerization for development and deployment
- **Make** - Build automation
- **Git Hooks** - Pre-commit linting
- **Bruno** - API testing and documentation

## Development Setup

### Local Development
1. Clone the repository
2. Copy `.env.example` to `.env` and configure database credentials
3. Run `go run main.go app:serve` to start the application

### Docker Development
1. Clone the repository
2. Configure `.env` file
3. Run `docker-compose up -d` to start the application

### Database Migrations
- Atlas handles schema migrations
- Migration commands available through Makefile:
  - `make migrate-status` - Check migration status
  - `make migrate-diff` - Generate migration file
  - `make migrate-apply` - Apply migrations
  - `make migrate-down` - Rollback migration
  - `make migrate-hash` - Hash migration files

## Technical Constraints

### Database
- MySQL is the primary database
- GORM for ORM operations
- Connection pooling configured in infrastructure/db.go
- Binary UUID used for primary keys

### API Design
- REST API with JSON responses
- Standardized response formats
- Proper error handling with custom error types
- Rate limiting support

### Authentication
- Cognito authentication integration
- JWT token validation
- Role-based access control

### File Storage
- AWS S3 integration for file storage
- Presigned URLs for secure file access

### Email
- AWS SES integration for email sending

## Dependencies

### Key Go Modules
- github.com/gin-gonic/gin - Web framework
- gorm.io/gorm - ORM
- gorm.io/driver/mysql - MySQL driver
- go.uber.org/fx - Dependency injection
- github.com/spf13/cobra - CLI commands
- github.com/spf13/viper - Configuration
- go.uber.org/zap - Logging
- github.com/aws/aws-sdk-go-v2 - AWS SDK
- github.com/lestrrat-go/jwx - JWT handling
- github.com/getsentry/sentry-go - Error monitoring
- github.com/joho/godotenv - Environment loading
- github.com/google/uuid - UUID generation
- github.com/ulule/limiter/v3 - Rate limiting

### External Services
- AWS S3 - File storage
- AWS Cognito - Authentication
- AWS SES - Email sending
- Sentry - Error tracking

## Environment Configuration

### Required Environment Variables
- `DB_USER` - Database username
- `DB_PASS` - Database password
- `DB_HOST` - Database host
- `DB_PORT` - Database port
- `DB_NAME` - Database name
- `SERVER_PORT` - Application server port
- `LOG_LEVEL` - Logging level
- `ENVIRONMENT` - Application environment (local, development, production)
- `SENTRY_DSN` - Sentry DSN for error tracking
- `STORAGE_BUCKET_NAME` - AWS S3 bucket name
- `AWS_REGION` - AWS region
- `AWS_ACCESS_KEY_ID` - AWS access key
- `AWS_SECRET_ACCESS_KEY` - AWS secret key

Additional environment variables are defined in the Env struct in pkg/framework/env.go.
