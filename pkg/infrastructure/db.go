package infrastructure

import (
	"clean-architecture/pkg/framework"
	"context"
	"fmt"
	"os"
	"time"

	"ariga.io/atlas-go-sdk/atlasexec"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database modal
type Database struct {
	*gorm.DB
	Logger framework.Logger
	Env    *framework.Env
}

// NewDatabase creates a new database instance
func NewDatabase(logger framework.Logger, env *framework.Env) Database {
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", env.DBUsername, env.DBPassword, env.DBHost, env.DBPort)

	logger.Info("opening db connection")
	db, err := gorm.Open(mysql.Open(url), &gorm.Config{Logger: logger.GetGormLogger()})
	if err != nil {
		logger.Panic(err)
	}

	logger.Info("creating database if it doesn't exist")
	if err = db.Exec("CREATE DATABASE IF NOT EXISTS " + env.DBName).Error; err != nil {
		logger.Info("couldn't create database")
		logger.Panic(err)
	}

	// close the current connection
	sqlDb, err := db.DB()
	if err != nil {
		logger.Panic(err)
	}
	if dbErr := sqlDb.Close(); dbErr != nil {
		logger.Panic(err)
	}

	// reopen connection with the given database, after creating or checking if the database exists
	logger.Info("using given database")
	urlWithDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", env.DBUsername, env.DBPassword, env.DBHost, env.DBPort, env.DBName)
	db, err = gorm.Open(mysql.Open(urlWithDB), &gorm.Config{Logger: logger.GetGormLogger()})
	if err != nil {
		logger.Panic(err)
	}

	conn, err := db.DB()
	if err != nil {
		logger.Info("couldn't get db connection")
		logger.Panic(err)
	}

	conn.SetConnMaxLifetime(time.Minute * 5)
	conn.SetMaxOpenConns(5)
	conn.SetMaxIdleConns(1)

	return Database{DB: db, Logger: logger, Env: env}
}

func NewMockDB() Database {
	_db, _, err := sqlmock.New()
	if err != nil {
		return Database{}
	}

	db, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      _db,
		SkipInitializeWithVersion: true,
	}))

	if err != nil {
		return Database{}
	}

	return Database{DB: db}
}

func (d *Database) RunMigration() {
	workdir, err := atlasexec.NewWorkingDir(
		atlasexec.WithMigrations(
			os.DirFS("./migrations"),
		),
	)
	if err != nil {
		d.Logger.Fatalf("failed to load working directory: %v", err)
	}
	// atlasexec works on a temporary directory, so we need to close it
	defer workdir.Close()

	// Initialize the client.
	client, err := atlasexec.NewClient(workdir.Path(), "atlas")
	if err != nil {
		d.Logger.Fatalf("failed to initialize client: %v", err)
	}

	res, err := client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
		URL:       fmt.Sprintf("mysql://%s:%s@%s:%s/%s?charset=utf8mb4&parseTime=True&loc=Local", d.Env.DBUsername, d.Env.DBPassword, d.Env.DBHost, d.Env.DBPort, d.Env.DBName),
		ExecOrder: "non-linear",
	})

	if err != nil {
		d.Logger.Fatalf("failed to apply migrations: %v", err)
	}

	d.Logger.Infof("Applied %d migrations\n", len(res.Applied))
}
