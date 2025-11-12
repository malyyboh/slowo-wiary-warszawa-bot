package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/database"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/models"
)

type UserRepository interface {
	AddOrUpdate(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userID int64) (*models.User, error)
	GetAll(ctx context.Context) ([]models.User, error)
	GetActive(ctx context.Context) ([]models.User, error)
	GetStats(ctx context.Context) (*models.UserStats, error)
	UpdateLastSeen(ctx context.Context, userID int64) error
	SetActive(ctx context.Context, userID int64, isActive bool) error
	SetBlocked(ctx context.Context, userID int64, isBlocked bool) error
	ExportDB(ctx context.Context) ([]byte, error)
}

type userRepository struct{}

func NewUserRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) AddOrUpdate(ctx context.Context, user *models.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (user_id, username, first_name, subscribed_at, is_active, is_blocked, last_seen, updated_at)
		VALUES (:user_id, :username, :first_name, :subscribed_at, :is_active, :is_blocked, :last_seen, :updated_at)
		ON CONFLICT(user_id) DO UPDATE SET
			username = :username,
			first_name = :first_name,
			last_seen = :last_seen,
			updated_at = :updated_at
	`

	_, err := database.DB.NamedExecContext(ctx, query, user)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("database write timeout for user %d: %w", user.UserID, err)
		}
		return fmt.Errorf("failed to add or update user %d: %w", user.UserID, err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userID int64) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var user models.User
	query := `SELECT * FROM users WHERE user_id = ?`

	err := database.DB.GetContext(ctx, &user, query, userID)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("database read timeout for user %d: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to get user %d: %w", userID, err)
	}

	return &user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var users []models.User
	query := `SELECT * FROM users ORDER BY subscribed_at DESC`

	err := database.DB.SelectContext(ctx, &users, query)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("database read timeout for all users: %w", err)
		}
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	return users, nil
}

func (r *userRepository) GetActive(ctx context.Context) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var users []models.User
	query := `
		SELECT * FROM users
		WHERE is_active = 1 AND is_blocked = 0
		ORDER BY subscribed_at DESC
	`

	err := database.DB.SelectContext(ctx, &users, query)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("database read timeout for active users: %w", err)
		}
		return nil, fmt.Errorf("failed to get active users: %w", err)
	}

	return users, nil
}

func (r *userRepository) GetStats(ctx context.Context) (*models.UserStats, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := `
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN is_active = 1 AND is_blocked = 0 THEN 1 ELSE 0 END) as active,
			SUM(CASE WHEN is_active = 0 AND is_blocked = 0 THEN 1 ELSE 0 END) as unsubscribed,
			SUM(CASE WHEN is_blocked = 1 THEN 1 ELSE 0 END) as blocked
		FROM users
		`
	stats := &models.UserStats{}
	err := database.DB.GetContext(ctx, stats, query)
	if err != nil {
		if err == context.DeadlineExceeded {
			return nil, fmt.Errorf("database read timeout for user stats: %w", err)
		}
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return stats, nil
}

func (r *userRepository) UpdateLastSeen(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	location, _ := time.LoadLocation("Europe/Warsaw")
	now := time.Now().In(location)

	query := `UPDATE users SET last_seen = ?, updated_at = ? WHERE user_id = ?`
	_, err := database.DB.ExecContext(ctx, query, now, now, userID)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("database write timeout updating last seen for user %d: %w", userID, err)
		}
		return fmt.Errorf("failed to update last seen for user %d: %w", userID, err)
	}
	return nil
}

func (r *userRepository) SetActive(ctx context.Context, userID int64, isActive bool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	location, _ := time.LoadLocation("Europe/Warsaw")
	now := time.Now().In(location)

	query := `UPDATE users SET is_active = ?, updated_at = ? WHERE user_id = ?`
	_, err := database.DB.ExecContext(ctx, query, isActive, now, userID)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("database write timeout setting active status for user %d: %w", userID, err)
		}
		return fmt.Errorf("failed to set active status for user %d: %w", userID, err)
	}
	return nil
}

func (r *userRepository) SetBlocked(ctx context.Context, userID int64, isBlocked bool) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	location, _ := time.LoadLocation("Europe/Warsaw")
	now := time.Now().In(location)

	query := `UPDATE users SET is_blocked = ?, updated_at = ? WHERE user_id = ?`
	_, err := database.DB.ExecContext(ctx, query, isBlocked, now, userID)
	if err != nil {
		if err == context.DeadlineExceeded {
			return fmt.Errorf("database write timeout setting blocked status for user %d: %w", userID, err)
		}
		return fmt.Errorf("failed to set blocked status for user %d: %w", userID, err)
	}
	return nil
}

func (r *userRepository) ExportDB(ctx context.Context) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var dbPath string
	err := database.DB.GetContext(ctx, &dbPath, "PRAGMA database_list")

	if err != nil || dbPath == "" {
		dbPath = os.Getenv("DATABASE_PATH")
		if dbPath == "" {
			dbPath = "./data/bot.db"
		}
	}

	log.Printf("🔍 ExportDB: Active DB connection path = '%s'", dbPath)

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Printf("❌ ExportDB: file not found at %s", dbPath)
		return nil, fmt.Errorf("database file not found: %s", dbPath)
	}

	info, _ := os.Stat(dbPath)
	log.Printf("🔍 ExportDB: file size = %d bytes", info.Size())

	data, err := os.ReadFile(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read database file: %w", err)
	}

	log.Printf("🔍 ExportDB: successfully read %d bytes", len(data))

	return data, nil
}
