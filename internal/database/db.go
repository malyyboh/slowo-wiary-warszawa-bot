package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func InitDB(dbPath string) error {
	var err error

	DB, err = sqlx.Connect("sqlite3", dbPath)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return err
	}

	log.Println("Database connected successfully")

	if err = createTables(); err != nil {
		return err
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
	`

	_, err := DB.Exec(schema)
	if err != nil {
		return err
	}

	log.Println("Database tables created successfully")
	return nil
}

func CloseDB() error {
	if DB != nil {
		log.Println("Closing database connection")
		return DB.Close()
	}
	return nil
}
