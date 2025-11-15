package models

import "time"

type RecurringEvent struct {
	ID          int    `db:"id" json:"id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`

	DayOfWeek int    `db:"day_of_week" json:"day_of_week"`
	EventTime string `db:"event_time" json:"event_time"`

	ReminderDayOffset int    `db:"reminder_day_offset" json:"reminder_day_offset"`
	ReminderTime      string `db:"reminder_time" json:"reminder_time"`

	Location        string `db:"location" json:"location"`
	Category        string `db:"category" json:"category"`
	RegistrationURL string `db:"registration_url" json:"registration_url"`

	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	CreatedBy int64     `db:"created_by" json:"created_by"`
}

// Helper functions

func (e *RecurringEvent) GetDayName() string {
	days := []string{"Неділя", "Понеділок", "Вівторок", "Середа", "Четвер", "П'ятниця", "Субота"}
	return days[e.DayOfWeek]
}

func (e *RecurringEvent) GetReminderDayName() string {
	reminderDay := e.DayOfWeek + e.ReminderDayOffset
	if reminderDay < 0 {
		reminderDay += 7
	}

	reminderDay = reminderDay % 7

	days := []string{"Неділя", "Понеділок", "Вівторок", "Середа", "Четвер", "П'ятниця", "Субота"}
	return days[reminderDay]
}
