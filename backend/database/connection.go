package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

// InitDB initializes the database connection
func InitDB() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://arvfinder:arvfinder_dev@localhost:5432/arvfinder?sslmode=disable"
	}

	var err error
	db, err = sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// RunMigrations runs database migrations
func RunMigrations(db *sql.DB) error {
	// Check if tables already exist (they should be created by Docker init script)
	var tableExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users')").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("failed to check for existing tables: %w", err)
	}

	if tableExists {
		log.Println("Database tables already exist, skipping migrations")
		return nil
	}

	// If tables don't exist, try to read and execute the schema file
	schemaPath := "./database/schema.sql"
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Printf("Schema file not found, assuming database is initialized by Docker: %v", err)
		return nil // Don't fail if schema file doesn't exist in container
	}

	// Execute the schema
	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}