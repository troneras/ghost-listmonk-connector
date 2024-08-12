package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

var db *sql.DB

func InitDB() error {
	config := utils.GetConfig()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure the connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

func GetDB() *sql.DB {
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

func RunMigrations(direction string, steps int) error {
    err := InitDB()
	if err != nil {
		utils.ErrorLogger.Fatalf("Failed to initialize database: %v", err)
	}
	defer GetDB().Close()
	utils.InfoLogger.Printf("Starting migration process. Direction: %s, Steps: %d", direction, steps)

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create migration driver: %v", err)
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"mysql",
		driver,
	)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create migration instance: %v", err)
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			utils.ErrorLogger.Printf("Error closing migration source: %v", srcErr)
		}
		if dbErr != nil {
			utils.ErrorLogger.Printf("Error closing migration database: %v", dbErr)
		}
	}()

	// Check if any migration files exist
	migrations, err := os.ReadDir("database/migrations")
	if err != nil {
		utils.ErrorLogger.Printf("Failed to read migrations directory: %v", err)
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	if len(migrations) == 0 {
		utils.InfoLogger.Println("No migration files found. Skipping migration process.")
		return nil
	}

	var migrationErr error
	switch direction {
	case "up":
		if steps > 0 {
			migrationErr = m.Steps(steps)
		} else {
			migrationErr = m.Up()
		}
	case "down":
		if steps > 0 {
			migrationErr = m.Steps(-steps)
		} else {
			migrationErr = m.Down()
		}
	default:
		return fmt.Errorf("invalid migration direction: %s", direction)
	}

	if migrationErr != nil && migrationErr != migrate.ErrNoChange {
		utils.ErrorLogger.Printf("Error applying migrations: %v", migrationErr)
		return fmt.Errorf("error applying migrations: %w", migrationErr)
	}

	utils.InfoLogger.Println("Migrations applied successfully")
	return nil
}
