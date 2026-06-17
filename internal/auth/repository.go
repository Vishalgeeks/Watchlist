package auth

import (
	"database/sql"
	"errors"
	"watchlist-backend/pkg/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(user *models.User, passwordHash string) error {
	query := `
		INSERT INTO users (name, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`
	return r.db.QueryRow(query, user.Name, user.Email, passwordHash).Scan(&user.ID)
}

func (r *Repository) GetUserByEmail(email string) (*models.User, string, error) {
	query := `
		SELECT id, name, email, password_hash, created_at
		FROM users WHERE email = $1
	`
	user := &models.User{}
	var passwordHash string

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Name, &user.Email, &passwordHash, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, "", errors.New("user not found")
	}
	if err != nil {
		return nil, "", err
	}
	return user, passwordHash, nil
}