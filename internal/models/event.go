package models

import "time"

type Event struct {
	ID              int       `db:"id"`
	Title           string    `db:"title"`
	Description     string    `db:"description"`
	Date            time.Time `db:"date"`
	Location        *string   `db:"location"`
	Category        *string   `db:"category"`
	RegistrationURL *string   `db:"registration_url"`
	IsPublished     bool      `db:"is_published"`
	CreatedAt       time.Time `db:"created_at"`
	CreatedBy       int64     `db:"created_by"`
}
