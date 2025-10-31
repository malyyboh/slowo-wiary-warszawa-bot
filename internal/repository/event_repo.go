package repository

import (
	"time"

	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/database"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/models"
)

type EventRepository interface {
	Create(event *models.Event) error
	GetByID(id int) (*models.Event, error)
	GetAll() ([]models.Event, error)
	GetUpcoming() ([]models.Event, error)
	Update(event *models.Event) error
	Delete(id int) error
}

type eventRepository struct{}

func NewEventRepository() EventRepository {
	return &eventRepository{}
}

func (r *eventRepository) Create(event *models.Event) error {
	query := `
		INSERT INTO events (title, description, date, location, category, registration_url, is_published, created_at, created_by)
		VALUES (:title, :description, :date, :location, :category, :registration_url, :is_published, :created_at, :created_by)
	`

	result, err := database.DB.NamedExec(query, event)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	event.ID = int(id)

	return nil
}

func (r *eventRepository) GetByID(id int) (*models.Event, error) {
	var event models.Event
	query := `SELECT * FROM events WHERE id = ?`

	err := database.DB.Get(&event, query, id)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *eventRepository) GetAll() ([]models.Event, error) {
	var events []models.Event
	query := `
		SELECT * FROM events 
		ORDER BY
			CASE WHEN date >= datetime('now') THEN 0 ELSE 1 END,
			date ASC
	`

	err := database.DB.Select(&events, query)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) GetUpcoming() ([]models.Event, error) {
	var events []models.Event
	now := time.Now().UTC()

	query := `
		SELECT * FROM events
		WHERE date >= ? AND is_published = 1
		ORDER BY date ASC
	`

	err := database.DB.Select(&events, query, now)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepository) Update(event *models.Event) error {
	query := `
		UPDATE events
		SET title = :title,
			description = :description,
			date = :date,
			location = :location,
			category = :category,
			registration_url = :registration_url,
			is_published = :is_published
		WHERE id = :id
	`
	_, err := database.DB.NamedExec(query, event)
	return err
}

func (r *eventRepository) Delete(id int) error {
	query := `DELETE FROM events WHERE id = ?`
	_, err := database.DB.Exec(query, id)
	return err
}
