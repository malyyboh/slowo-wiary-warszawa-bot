package models

import "time"

type User struct {
	UserID       int64     `db:"user_id"`
	Username     string    `db:"username"`
	FirstName    string    `db:"first_name"`
	SubscribedAt time.Time `db:"subscribed_at"`
	IsActive     bool      `db:"is_active"`
	IsBlocked    bool      `db:"is_blocked"`
	LastSeen     time.Time `db:"last_seen"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UserStats struct {
	Total        int
	Active       int
	Unsubscribed int
	Blocked      int
}
