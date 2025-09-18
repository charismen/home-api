package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// DB is the global database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set connection pool parameters
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to MySQL database")

	// Create tables if they don't exist
	createTables()
}

// createTables creates the necessary tables if they don't exist
func createTables() {
	// Create items table
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS items (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(100) NOT NULL,
		external_id VARCHAR(100) NOT NULL UNIQUE,
		data TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		INDEX idx_external_id (external_id)
	)
	`)
	if err != nil {
		log.Fatalf("Failed to create items table: %v", err)
	}

	// Create orders table for Part 3
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS orders (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		customer_id VARCHAR(36) NOT NULL,
		amount DECIMAL(10,2) NOT NULL,
		status ENUM('PENDING', 'PAID', 'CANCELLED') NOT NULL,
		created_at DATETIME NOT NULL
	)
	`)
	if err != nil {
		log.Fatalf("Failed to create orders table: %v", err)
	}
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}