package pgstorage

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ibeloyar/gophkeeper/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetSecret_FoundText(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}

	title := "secret1"
	userID := int64(1)
	ctx := context.Background()
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "title", "metadata", "created_at", "updated_at",
		"secret_type", "login", "password", "text_data", "binary_data",
		"card_number", "card_exp", "card_holder",
	}).AddRow(int64(123), userID, title, "meta", createdAt, createdAt, "text", "", "", "text data", nil, "", "", "")

	mock.ExpectQuery(`SELECT \* from secrets WHERE title = \$1 AND user_id = \$2`).
		WithArgs(title, userID).
		WillReturnRows(rows)

	result, err := repo.GetSecret(ctx, title, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(123), result.ID)
	assert.Equal(t, "text data", result.TextData)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetSecret_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}

	ctx := context.Background()

	mock.ExpectQuery(`SELECT \* from secrets WHERE title = \$1 AND user_id = \$2`).
		WithArgs("secret1", int64(1)).
		WillReturnError(sql.ErrNoRows)

	result, err := repo.GetSecret(ctx, "secret1", 1)

	assert.NoError(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetSecret_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}

	ctx := context.Background()

	mock.ExpectQuery(`SELECT \* from secrets WHERE title = \$1 AND user_id = \$2`).
		WithArgs("secret1", int64(1)).
		WillReturnError(errors.New("connection failed"))

	result, err := repo.GetSecret(ctx, "secret1", 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetSecret_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}

	ctx := context.Background()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "title", "metadata", "created_at", "updated_at",
		"secret_type", "login", "password", "text_data", "binary_data",
		"card_number", "card_exp", "card_holder",
	}).AddRow("invalid_id", 1, "title", "", time.Now(), time.Now(), 1, "", "", "", nil, "", "", "")

	mock.ExpectQuery(`SELECT \* from secrets WHERE title = \$1 AND user_id = \$2`).
		WithArgs("secret1", int64(1)).
		WillReturnRows(rows)

	result, err := repo.GetSecret(ctx, "secret1", 1)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetSecrets_Empty(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}
	ctx := context.Background()

	mock.ExpectQuery(`SELECT id, user_id, title, metadata, created_at, updated_at, secret_type, login, password_hash, text_data, binary_data, card_number, card_exp, card_holder FROM secrets WHERE user_id = \$1 ORDER BY title`).
		WithArgs(int64(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "title", "metadata", "created_at", "updated_at",
			"secret_type", "login", "password_hash", "text_data", "binary_data",
			"card_number", "card_exp", "card_holder",
		}))

	secrets, err := repo.GetSecrets(ctx, 1)

	assert.NoError(t, err)
	assert.Empty(t, secrets)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetSecrets_FoundText(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}
	ctx := context.Background()
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "title", "metadata", "created_at", "updated_at",
		"secret_type", "login", "password_hash", "text_data", "binary_data",
		"card_number", "card_exp", "card_holder",
	}).AddRow(123, int64(1), "text secret", "meta", createdAt, createdAt, "text", "", "", "data", nil, "", "", "")

	mock.ExpectQuery(`SELECT id, user_id, title, metadata, created_at, updated_at, secret_type, login, password_hash, text_data, binary_data, card_number, card_exp, card_holder FROM secrets WHERE user_id = \$1 ORDER BY title`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	secrets, err := repo.GetSecrets(ctx, 1)

	assert.NoError(t, err)
	assert.Len(t, secrets, 1)
	assert.Equal(t, int64(123), secrets[0].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetSecrets_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}
	ctx := context.Background()

	mock.ExpectQuery(`SELECT id, user_id, title, metadata, created_at, updated_at, secret_type, login, password_hash, text_data, binary_data, card_number, card_exp, card_holder FROM secrets WHERE user_id = \$1 ORDER BY title`).
		WithArgs(int64(1)).
		WillReturnError(errors.New("db error"))

	secrets, err := repo.GetSecrets(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, secrets)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_GetSecrets_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{
		"id", "user_id", "title", "metadata", "created_at", "updated_at",
		"secret_type", "login", "password_hash", "text_data", "binary_data",
		"card_number", "card_exp", "card_holder",
	}).AddRow("invalid", 1, "title", "", time.Now(), time.Now(), "text", "", "", "", nil, "", "", "")

	mock.ExpectQuery(`SELECT id, user_id, title, metadata, created_at, updated_at, secret_type, login, password_hash, text_data, binary_data, card_number, card_exp, card_holder FROM secrets WHERE user_id = \$1 ORDER BY title`).
		WithArgs(int64(1)).
		WillReturnRows(rows)

	secrets, err := repo.GetSecrets(ctx, 1)

	assert.Error(t, err)
	assert.Nil(t, secrets)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateSecret_SuccessText(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}
	ctx := context.Background()

	secret := &model.CreateSecretDTO{
		UserID:     1,
		Title:      "text secret",
		SecretType: model.SecretTypeText,
		TextData:   "data",
		Metadata:   "meta",
	}

	query := `INSERT INTO secrets (user_id,title,metadata,secret_type,login,password_hash, text_data,binary_data,card_number,card_exp,card_holder ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(
			secret.UserID,
			secret.Title,
			secret.Metadata,
			secret.SecretType,
			"", // login
			"", // password_hash
			secret.TextData,
			[]byte(nil), // binary_data
			"",          // card_number
			"",          // card_exp
			"",          // card_holder
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(123))

	id, err := repo.CreateSecret(ctx, secret)

	assert.NoError(t, err)
	assert.Equal(t, int64(123), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_CreateSecret_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}
	ctx := context.Background()

	secret := &model.CreateSecretDTO{
		UserID:     1,
		Title:      "text secret",
		SecretType: model.SecretTypeText,
		TextData:   "data",
		Metadata:   "meta",
	}

	query := `INSERT INTO secrets (user_id,title,metadata,secret_type,login,password_hash, text_data,binary_data,card_number,card_exp,card_holder ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;`

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(
			secret.UserID,
			secret.Title,
			secret.Metadata,
			secret.SecretType,
			"", // login
			"", // password_hash
			secret.TextData,
			[]byte(nil), // binary_data
			"",          // card_number
			"",          // card_exp
			"",          // card_holder
		).
		WillReturnError(errors.New("database insert failed"))

	id, err := repo.CreateSecret(ctx, secret)

	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.Contains(t, err.Error(), "database insert failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteSecret_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}
	ctx := context.Background()

	title := "my-secret"
	userID := int64(123)

	q := `DELETE FROM secrets WHERE user_id = $1 AND title = $2;`

	rows := sqlmock.NewResult(0, 1) // 1 строка удалена
	mock.ExpectExec(regexp.QuoteMeta(q)).
		WithArgs(userID, title).
		WillReturnResult(rows)

	err = repo.DeleteSecret(ctx, title, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRepository_DeleteSecret_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}
	ctx := context.Background()

	title := "nonexistent-secret"
	userID := int64(123)

	q := `DELETE FROM secrets WHERE user_id = $1 AND title = $2;`

	rows := sqlmock.NewResult(0, 0) // 0 строк удалено
	mock.ExpectExec(regexp.QuoteMeta(q)).
		WithArgs(userID, title).
		WillReturnResult(rows)

	err = repo.DeleteSecret(ctx, title, userID)

	assert.Error(t, err)
}

func TestRepository_DeleteSecret_DBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &PGStorage{db: db}
	ctx := context.Background()

	title := "my-secret"
	userID := int64(123)

	q := `DELETE FROM secrets WHERE user_id = $1 AND title = $2;`

	mock.ExpectExec(regexp.QuoteMeta(q)).
		WithArgs(userID, title).
		WillReturnError(errors.New("database error"))

	err = repo.DeleteSecret(ctx, title, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "database error")
	assert.NoError(t, mock.ExpectationsWereMet())
}
