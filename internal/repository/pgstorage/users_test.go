package pgstorage

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetUserByLogin_Found(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}

	createdAt := time.Now()
	rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "created_at", "updated_at"}).
		AddRow(123, "testuser", "hashedpass", createdAt, createdAt)

	mock.ExpectQuery(`SELECT \* FROM users WHERE login = \$1`).
		WithArgs("testuser").
		WillReturnRows(rows)

	result := repo.GetUserByLogin(context.Background(), "testuser")

	assert.NotNil(t, result)
	assert.Equal(t, int64(123), result.ID)
	assert.Equal(t, "testuser", result.Login)
	assert.Equal(t, "hashedpass", result.PasswordHash)
	assert.WithinDuration(t, result.CreatedAt, createdAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserByLogin_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}

	mock.ExpectQuery(`SELECT \* FROM users WHERE login = \$1`).
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	result := repo.GetUserByLogin(context.Background(), "nonexistent")

	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserByLogin_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}

	mock.ExpectQuery(`SELECT \* FROM users WHERE login = \$1`).
		WithArgs("testuser").
		WillReturnError(errors.New("connection failed"))

	result := repo.GetUserByLogin(context.Background(), "testuser")

	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetUserByLogin_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}

	rows := sqlmock.NewRows([]string{"id", "login", "password_hash", "created_at", "updated_at"}).
		AddRow("invalid_id", "testuser", "hashed", time.Now(), time.Now())

	mock.ExpectQuery(`SELECT \* FROM users WHERE login = \$1`).
		WithArgs("testuser").
		WillReturnRows(rows)

	result := repo.GetUserByLogin(context.Background(), "testuser")

	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}

	user := model.User{
		Login:        "testuser",
		PasswordHash: "hashedpass",
	}
	ctx := context.Background()

	mock.ExpectQuery(`INSERT INTO users \(login, password_hash\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs(user.Login, user.PasswordHash).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

	id, err := repo.CreateUser(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, int64(123), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateUser_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}

	user := model.User{Login: "testuser", PasswordHash: "hashed"}
	ctx := context.Background()

	mock.ExpectQuery(`INSERT INTO users \(login, password_hash\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs(user.Login, user.PasswordHash).
		WillReturnError(errors.New("duplicate key"))

	id, err := repo.CreateUser(ctx, user)

	assert.Error(t, err)
	assert.Zero(t, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateUser_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}

	user := model.User{Login: "testuser", PasswordHash: "hashed"}
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id"}).AddRow("invalid_id")
	mock.ExpectQuery(`INSERT INTO users \(login, password_hash\) VALUES \(\$1, \$2\) RETURNING id`).
		WithArgs(user.Login, user.PasswordHash).
		WillReturnRows(rows)

	id, err := repo.CreateUser(ctx, user)

	assert.Error(t, err)
	assert.Zero(t, id)
	assert.NoError(t, mock.ExpectationsWereMet())
}
