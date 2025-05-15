package testutil

import (
	"clean-architecture/pkg/framework"
	"clean-architecture/pkg/infrastructure"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type MySQLContainer struct {
	Container tc.Container
	Host      string
	Port      string
	Cleanup   func()
}

func SetupMySQLContainer(ctx context.Context, env *framework.Env) (*MySQLContainer, error) {
	log.Printf("Setting up MySQL test container with image: %s", env.DBName)

	req := tc.ContainerRequest{
		Image: "mysql/mysql-server:8.0",
		Env: map[string]string{
			"MYSQL_DATABASE":      env.DBName,
			"MYSQL_USER":          env.DBUsername,
			"MYSQL_PASSWORD":      env.DBPassword,
			"MYSQL_ROOT_PASSWORD": env.DBPassword,
		},
		ExposedPorts: []string{"3306/tcp"},
		WaitingFor: wait.ForLog("ready for connections").
			WithOccurrence(1).
			WithStartupTimeout(120 * time.Second),
		// FromDockerfile: tc.FromDockerfile{
		// 	Context:    "./docker",
		// 	Dockerfile: "db.Dockerfile",
		// },
	}

	mysqlContainer, err := tc.GenericContainer(ctx, tc.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start container: %w", err)
	}

	cleanup := func() {
		log.Println("Cleaning up MySQL container...")
		if err := mysqlContainer.Terminate(ctx); err != nil {
			log.Printf("Error terminating container: %v", err)
		} else {
			log.Println("Container terminated.")
		}
	}

	host, err := mysqlContainer.Host(ctx)
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := mysqlContainer.MappedPort(ctx, "3306")
	if err != nil {
		cleanup()
		return nil, fmt.Errorf("failed to get container port: %w", err)
	}

	return &MySQLContainer{
		Container: mysqlContainer,
		Host:      host,
		Port:      port.Port(),
		Cleanup:   cleanup,
	}, nil
}

// ConnectToDatabase establishes a connection to the database
func ConnectToDatabase(
	ctx context.Context,
	container *MySQLContainer,
	env *framework.Env,
) (*gorm.DB, error) {
	log.Printf("Container host: %s, port: %s", container.Host, container.Port)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		env.DBUsername,
		env.DBPassword,
		net.JoinHostPort(container.Host, container.Port),
		env.DBName,
	)
	log.Printf("Attempting to connect to DSN: tcp(%s)/%s", net.JoinHostPort(container.Host, container.Port), env.DBName)

	const maxRetries = 10
	const retryDelay = 5 * time.Second

	for i := range maxRetries {
		db, openErr := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if openErr == nil {
			sqlDB, sqlErr := db.DB()
			if sqlErr == nil {
				if pingErr := sqlDB.PingContext(ctx); pingErr == nil {
					log.Println("Successfully connected to test database.")
					return db, nil
				}
			}
			// Close temporary failed connection
			if sqlDB != nil {
				sqlDB.Close()
			}
		}

		log.Printf("Database connection attempt %d failed. Retrying in %s...", i+1, retryDelay)
		time.Sleep(retryDelay)

		if i == maxRetries-1 {
			return nil, fmt.Errorf("failed to connect to database after %d attempts", maxRetries)
		}
	}

	return nil, fmt.Errorf("unexpected error in connection loop")
}

// NewTestDatabase creates a test database using a MySQL container
func NewTestDatabase(
	logger framework.Logger,
	env *framework.Env,
) infrastructure.Database {
	logger.Info("Creating test database...")
	ctx := context.Background()
	container, err := SetupMySQLContainer(ctx, env)
	if err != nil {
		log.Printf("Failed to setup MySQL container: %v", err)
		return infrastructure.Database{}
	}

	gormDB, err := ConnectToDatabase(ctx, container, env)
	if err != nil {
		container.Cleanup()
		log.Printf("Failed to connect to database: %v", err)
		return infrastructure.Database{}
	}
	logger.Info("Connected to test database.")

	// Get the underlying *sql.DB and combine its closure with container cleanup
	sqlDB, err := gormDB.DB()
	if err != nil {
		container.Cleanup()
		return infrastructure.Database{}
	}

	// Set a longer connection timeout
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)

	originalCleanup := container.Cleanup
	container.Cleanup = func() {
		log.Println("Closing GORM database connection...")
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("GORM database connection closed.")
		}
		originalCleanup()
	}

	env.DBHost = container.Host
	env.DBPort = container.Port
	return infrastructure.Database{
		DB:     gormDB,
		Logger: logger,
		Env:    env,
	}
}
