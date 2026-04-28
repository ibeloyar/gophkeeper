package pgstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ibeloyar/gophkeeper/internal/model"
)

// CreateUser inserts new user into PostgreSQL users table and returns generated ID.
// Uses RETURNING clause for efficient single-query ID retrieval. Expects password_hash
// to be pre-hashed (bcrypt). Returns PostgreSQL error on constraint violations.
func (s *PGStorage) CreateUser(ctx context.Context, user model.User) (int64, error) {
	var userID int64

	err := s.db.QueryRowContext(ctx, `INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`,
		user.Login, user.PasswordHash,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

// GetUserByLogin retrieves user by exact login match. Returns nil if user doesn't exist
// (sql.ErrNoRows) or on scan errors. Scans all user fields including timestamps.
// Designed for authentication lookup - returns pointer to avoid nil panics.
func (s *PGStorage) GetUserByLogin(ctx context.Context, login string) *model.User {
	var user model.User

	row := s.db.QueryRowContext(ctx, "SELECT * FROM users WHERE login = $1", login)

	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return nil
	}

	return &user
}
