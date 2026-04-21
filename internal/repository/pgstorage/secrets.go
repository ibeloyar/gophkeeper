package pgstorage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ibeloyar/gophkeeper/internal/model"
)

func (s *PGStorage) GetSecret(ctx context.Context, title string, userID int64) (*model.Secret, error) {
	var secret model.Secret

	row := s.db.QueryRowContext(ctx, "SELECT * from secrets WHERE title = $1 AND user_id = $2;", title, userID)

	err := row.Scan(&secret.ID, &secret.UserID, &secret.Title, &secret.Metadata, &secret.CreatedAt, &secret.UpdatedAt,
		&secret.SecretType, &secret.Login, &secret.Password, &secret.TextData, &secret.BinaryData, &secret.CardNumber,
		&secret.CardExp, &secret.CardHolder)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &secret, nil
}

func (s *PGStorage) GetSecrets(ctx context.Context, userID int64) ([]*model.Secret, error) {
	var secrets []*model.Secret

	rows, err := s.db.QueryContext(ctx, `SELECT
       id, user_id, title, metadata, created_at, updated_at,
       secret_type, login, password_hash, text_data,
       binary_data, card_number, card_exp, card_holder
   FROM secrets
   WHERE user_id = $1
   ORDER BY title;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var secret model.Secret
		err := rows.Scan(
			&secret.ID, &secret.UserID, &secret.Title, &secret.Metadata,
			&secret.CreatedAt, &secret.UpdatedAt,
			&secret.SecretType,
			&secret.Login, &secret.Password,
			&secret.TextData,
			&secret.BinaryData,
			&secret.CardNumber, &secret.CardExp, &secret.CardHolder,
		)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, &secret)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return secrets, nil
}

func (s *PGStorage) CreateSecret(ctx context.Context, secret *model.CreateSecretDTO) (int64, error) {
	var id int64

	q := `INSERT INTO secrets (user_id,title,metadata,secret_type,login,password_hash,
			text_data,binary_data,card_number,card_exp,card_holder
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;`

	if err := s.db.QueryRowContext(ctx, q, secret.UserID, secret.Title, secret.Metadata, string(secret.SecretType),
		secret.Login, secret.Password,
		secret.TextData,
		secret.BinaryData,
		secret.CardNumber, secret.CardExp, secret.CardHolder,
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (s *PGStorage) DeleteSecret(ctx context.Context, title string, userID int64) error {
	q := `DELETE FROM secrets WHERE user_id = $1 AND title = $2;`

	res, err := s.db.ExecContext(ctx, q, userID, title)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrSecretNotFound
	}

	return nil
}
