package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func InitDB(dbPath string) error {

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %w", err)
	}
	var err error

	dsn := dbPath + "?_busy_timeout=5000&_journal_mode=WAL&_sync=NORMAL"

	DB, err = sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB.SetMaxOpenConns(1)
	DB.SetMaxIdleConns(1)

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully")

	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	return nil
}

func createTables() error {
	schema := `
	CREATE TABLE IF NOT EXISTS events (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		date DATETIME NOT NULL,
		location TEXT,
		category TEXT,
		registration_url TEXT,
		is_published BOOLEAN NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		created_by INTEGER NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_events_date ON events(date);
	CREATE INDEX IF NOT EXISTS idx_events_is_published ON events(is_published);
	
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY,
		username TEXT,
		first_name TEXT NOT NULL,
		subscribed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		is_active BOOLEAN NOT NULL DEFAULT 1,
		is_blocked BOOLEAN NOT NULL DEFAULT 0,
		last_seen DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
	CREATE INDEX IF NOT EXISTS idx_users_is_blocked ON users(is_blocked);
	
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	log.Println("Database tables created successfully")

	if err := runMigrations(); err != nil {
		return err
	}

	return nil
}

func runMigrations() error {
	migrations := []struct {
		version int
		sql     string
	}{
		// {1, "ALTER TABLE users ADD COLUMN phone TEXT;"},
		// Add new migrations here in the future
	}

	for _, m := range migrations {
		var exists int
		err := DB.Get(&exists, "SELECT 1 FROM schema_migrations WHERE version = ?", m.version)
		if err == nil && exists == 1 {
			continue
		}

		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to check migration %d: %w", m.version, err)
		}

		_, err = DB.Exec(m.sql)
		if err != nil {
			return fmt.Errorf("migration %d failed: %w", m.version, err)
		}

		_, err = DB.Exec("INSERT INTO schema_migrations (version) VALUES (?)", m.version)
		if err != nil {
			return fmt.Errorf("failed to record migration %d: %w", m.version, err)
		}

		log.Printf("Migration %d applied successfully", m.version)
	}

	return nil
}

func CloseDB() error {
	if DB != nil {
		log.Println("Closing database connection")
		return DB.Close()
	}
	return nil
}
