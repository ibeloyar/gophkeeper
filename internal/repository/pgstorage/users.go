package pgstorage

import (
	"database/sql"
	"errors"

	"github.com/ibeloyar/gophkeeper/internal/model"
)

func (s *PGStorage) CreateUser(user model.User) (int64, error) {
	var userID int64

	err := s.db.QueryRow(`INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id`,
		user.Login, user.PasswordHash,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *PGStorage) GetUserByLogin(login string) *model.User {
	var user model.User

	row := s.db.QueryRow("SELECT * FROM users WHERE login = $1", login)

	err := row.Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return nil
	}

	return &user
}
