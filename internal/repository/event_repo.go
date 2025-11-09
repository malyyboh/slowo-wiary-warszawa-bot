package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/database"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/models"
)

type EventRepository interface {
	Create(ctx context.Context, event *models.Event) error
	GetByID(ctx context.Context, id int) (*models.Event, error)
	GetAll(ctx context.Context) ([]models.Event, error)
	GetUpcoming(ctx context.Context) ([]models.Event, error)
	Update(ctx context.Context, event *models.Event) error
	Delete(ctx context.Context, id int) error
}

type eventRepository struct{}

func NewEventRepository() EventRepository {
	return &eventRepository{}
}

func (r *eventRepository) Create(ctx context.Context, event *models.Event) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO events (title, description, date, location, category, registration_url, is_published, created_at, created_by)
		VALUES (:title, :description, :date, :location, :category, :registration_url, :is_published, :created_at, :created_by)
	`

	result, err := database.DB.NamedExecContext(ctx, query, event)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("database write timeout creating event: %w", err)
		}
		return fmt.Errorf("failed to create event: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}
	event.ID = int(id)

	return nil
}

func (r *eventRepository) GetByID(ctx context.Context, id int) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var event models.Event
	query := `SELECT * FROM events WHERE id = ?`

	err := database.DB.GetContext(ctx, &event, query, id)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("database read timeout for event %d: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get event %d: %w", id, err)
	}

	return &event, nil
}

func (r *eventRepository) GetAll(ctx context.Context) ([]models.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var events []models.Event
	query := `
		SELECT * FROM events 
		ORDER BY
			CASE WHEN date >= datetime('now') THEN 0 ELSE 1 END,
			date ASC
	`

	err := database.DB.SelectContext(ctx, &events, query)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("database read timeout for all events: %w", err)
		}
		return nil, fmt.Errorf("failed to get all events: %w", err)
	}

	return events, nil
}

func (r *eventRepository) GetUpcoming(ctx context.Context) ([]models.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var events []models.Event

	location, _ := time.LoadLocation("Europe/Warsaw")
	now := time.Now().In(location)

	query := `
		SELECT * FROM events
		WHERE date >= ? AND is_published = 1
		ORDER BY date ASC
	`

	err := database.DB.SelectContext(ctx, &events, query, now)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("database read timeout for upcoming events: %w", err)
		}
		return nil, fmt.Errorf("failed to get upcoming events: %w", err)
	}

	return events, nil
}

func (r *eventRepository) Update(ctx context.Context, event *models.Event) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

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
	_, err := database.DB.NamedExecContext(ctx, query, event)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("database write timeout updating event %d: %w", event.ID, err)
		}
		return fmt.Errorf("failed to update event %d: %w", event.ID, err)
	}
	return nil
}

func (r *eventRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM events WHERE id = ?`

	_, err := database.DB.ExecContext(ctx, query, id)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("database write timeout deleting event %d: %w", id, err)
		}
		return fmt.Errorf("failed to delete event %d: %w", id, err)
	}
	return nil
}
