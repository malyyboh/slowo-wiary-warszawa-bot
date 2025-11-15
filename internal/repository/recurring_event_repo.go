package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/database"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/models"
)

type RecurringEventRepository interface {
	Create(ctx context.Context, event *models.RecurringEvent) error
	GetByID(ctx context.Context, id int) (*models.RecurringEvent, error)
	GetAll(ctx context.Context) ([]models.RecurringEvent, error)
	GetActive(ctx context.Context) ([]models.RecurringEvent, error)
	Update(ctx context.Context, event *models.RecurringEvent) error
	Delete(ctx context.Context, id int) error
	SetActive(ctx context.Context, id int, isActive bool) error
}

type recurringEventRepository struct{}

func NewRecurringEventRepository() RecurringEventRepository {
	return &recurringEventRepository{}
}

func (r *recurringEventRepository) Create(ctx context.Context, event *models.RecurringEvent) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO recurring_events (
			title, description, day_of_week, event_time,
			reminder_day_offset, reminder_time,
			location, category, registration_url,
			is_active, created_at, created_by
		)
		VALUES (
			:title, :description, :day_of_week, :event_time,
			:reminder_day_offset, :reminder_time,
			:location, :category, :registration_url,
			:is_active, :created_at, :created_by
		)
	`
	result, err := database.DB.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmt.Errorf("failed to create recurring event: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert id: %w", err)
	}

	event.ID = int(id)
	log.Printf("✅ Created recurring event: %s (ID: %d)", event.Title, event.ID)

	return nil
}

func (r *recurringEventRepository) GetByID(ctx context.Context, id int) (*models.RecurringEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var event models.RecurringEvent
	query := `SELECT * FROM recurring_events WHERE id = ?`

	err := database.DB.GetContext(ctx, &event, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get recurring event %d: %w", id, err)
	}

	return &event, nil
}

func (r *recurringEventRepository) GetAll(ctx context.Context) ([]models.RecurringEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var events []models.RecurringEvent
	query := `SELECT * FROM recurring_events ORDER BY day_of_week, event_time`

	err := database.DB.SelectContext(ctx, &events, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all recurring events: %w", err)
	}

	return events, nil
}

func (r *recurringEventRepository) GetActive(ctx context.Context) ([]models.RecurringEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var events []models.RecurringEvent
	query := `SELECT * FROM recurring_events WHERE is_active = 1 ORDER BY day_of_week, event_time`

	err := database.DB.SelectContext(ctx, &events, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active recurring events: %w", err)
	}

	return events, nil
}

func (r *recurringEventRepository) Update(ctx context.Context, event *models.RecurringEvent) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		UPDATE recurring_events SET
			title = :title,
			description = :description,
			day_of_week = :day_of_week,
			event_time = :event_time,
			reminder_day_offset = :reminder_day_offset,
			reminder_time = :reminder_time,
			location = :location,
			category = :category,
			registration_url = :registration_url,
			is_active = :is_active
		WHERE id = :id
	`
	_, err := database.DB.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmt.Errorf("failed to update recurring event %d: %w", event.ID, err)
	}

	log.Printf("✅ Updated recurring event: %s (ID: %d)", event.Title, event.ID)

	return nil
}

func (r *recurringEventRepository) Delete(ctx context.Context, id int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `DELETE FROM recurring_events WHERE id = ?`
	_, err := database.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete recurring event %d: %w", id, err)
	}

	log.Printf("✅ Deleted recurring event ID: %d", id)

	return nil
}

func (r *recurringEventRepository) SetActive(ctx context.Context, id int, isActive bool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `UPDATE recurring_events SET is_active = ? WHERE id = ?`
	_, err := database.DB.ExecContext(ctx, query, isActive, id)
	if err != nil {
		return fmt.Errorf("failed to set active status for recurring event %d: %w", id, err)
	}

	return nil
}
