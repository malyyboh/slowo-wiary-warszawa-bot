package repository

import (
	"time"

	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/database"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/models"
)

type UserRepository interface {
	AddOrUpdate(user *models.User) error
	GetByID(userID int64) (*models.User, error)
	GetAll() ([]models.User, error)
	GetActive() ([]models.User, error)
	GetStats() (*models.UserStats, error)
	UpdateLastSeen(userID int64) error
	SetActive(userID int64, isActive bool) error
	SetBlocked(userID int64, isBlocked bool) error
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) AddOrUpdate(user *models.User) error {
	query := `
		INSERT INTO users (user_id, username, first_name, subscribed_at, is_active, is_blocked, last_seen, updated_at)
		VALUES (:user_id, :username, :first_name, :subscribed_at, :is_active, :is_blocked, :last_seen, :updated_at)
		ON CONFLICT(user_id) DO UPDATE SET
			username = :username,
			first_name = :first_name,
			last_seen = :last_seen,
			updated_at = :updated_at
	`

	_, err := database.DB.NamedExec(query, user)
	return err
}

func (r *userRepository) GetByID(userID int64) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE user_id = ?`

	err := database.DB.Get(&user, query, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users ORDER BY subscribed_at DESC`

	err := database.DB.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetActive() ([]models.User, error) {
	var users []models.User
	query := `
		SELECT * FROM users
		WHERE is_active = 1 AND is_blocked = 0
		ORDER BY subscribed_at DESC
	`

	err := database.DB.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetStats() (*models.UserStats, error) {
	stats := &models.UserStats{}

	err := database.DB.Get(&stats.Total, `SELECT COUNT(*) FROM users`)
	if err != nil {
		return nil, err
	}

	err = database.DB.Get(&stats.Active, `SELECT COUNT(*) FROM users WHERE is_active = 1 AND is_blocked = 0`)
	if err != nil {
		return nil, err
	}

	err = database.DB.Get(&stats.Unsubscribed, `SELECT COUNT(*) FROM users WHERE is_active = 0 AND is_blocked = 0`)
	if err != nil {
		return nil, err
	}

	err = database.DB.Get(&stats.Blocked, `SELECT COUNT(*) FROM users WHERE is_blocked = 1`)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (r *userRepository) UpdateLastSeen(userID int64) error {
	query := `UPDATE users SET last_seen = ?, updated_at = ? WHERE user_id = ?`
	now := time.Now()
	_, err := database.DB.Exec(query, now, now, userID)
	return err
}

func (r *userRepository) SetActive(userID int64, isActive bool) error {
	query := `UPDATE users SET is_active = ?, updated_at = ? WHERE user_id = ?`
	_, err := database.DB.Exec(query, isActive, time.Now(), userID)
	return err
}

func (r *userRepository) SetBlocked(userID int64, isBlocked bool) error {
	query := `UPDATE users SET is_blocked = ?, updated_at = ? WHERE user_id = ?`
	_, err := database.DB.Exec(query, isBlocked, time.Now(), userID)
	return err
}
