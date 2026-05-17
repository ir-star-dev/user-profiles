package auth_postgres

import (
	"user-profiles/internal/auth"
	"user-profiles/internal/models"

	"github.com/jmoiron/sqlx"
)

type tokenRepository struct {
	db auth.DBTX
}

func NewTokenRepository(db *sqlx.DB) auth.RefreshRepository {
	return &tokenRepository{db: db}
}

func (repo *tokenRepository) Save(token *models.RefreshToken) error {
	query := `
        INSERT INTO tokens (token_hash, revoked, expires_at, user_id, family_id) 
        VALUES (:token_hash, :revoked, :expires_at, :user_id, :family_id)
    `
	_, err := repo.db.NamedExec(query, token)
	if err != nil {
		return err
	}
	return nil
}

func (repo *tokenRepository) WithTx(fn func(repo auth.RefreshRepository) error) error {
	tx, err := repo.db.(*sqlx.DB).Beginx()
	if err != nil {
		return err
	}

	txRepo := &tokenRepository{
		db: tx,
	}

	defer tx.Rollback()

	if err := fn(txRepo); err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *tokenRepository) FindTokenByHash(hash []byte) (*models.RefreshToken, error) {
	var token models.RefreshToken
	query := `SELECT id, expires_at, revoked, user_id, family_id FROM tokens WHERE token_hash = $1`
	err := repo.db.Get(&token, query, hash)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (repo *tokenRepository) FindUserIdByHash(hash []byte) (int, error) {
	var token models.RefreshToken
	query := `SELECT user_id FROM tokens WHERE token_hash = $1`
	err := repo.db.Get(&token, query, hash)
	if err != nil {
		return 0, err
	}
	return token.UserId, nil
}

func (repo *tokenRepository) Revoke(hash []byte, uId int) error {
	t := &models.RefreshToken{
		UserId:    uId,
		TokenHash: hash,
	}
	query := `
        UPDATE tokens 
        SET revoked = true
        WHERE user_id = :user_id AND token_hash = :token_hash
    `
	_, err := repo.db.NamedExec(query, t)
	if err != nil {
		return err
	}
	return nil
}

func (repo *tokenRepository) RevokeFamily(familyId string) error {
	query := `
		UPDATE tokens
		SET revoked = true
		WHERE family_id = $1
	`
	_, err := repo.db.Exec(query, familyId)
	return err
}