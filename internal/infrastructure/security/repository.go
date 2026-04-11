package security

import (
	"user-profiles/internal/domain/token"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	NamedExec(query string, arg any) (sql.Result, error)
	Get(dest any, query string, args ...any) error
}

type tokenRepository struct {
	db DBTX
}

func NewTokenRepository(db *sqlx.DB) token.Repository {
	return &tokenRepository{
		db: db,
	}
}

func (repo *tokenRepository) Save(token *token.RefreshToken) error {
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

func (repo *tokenRepository) WithTx(fn func(repo token.Repository) error) error {
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

func (repo *tokenRepository) FindTokenByHash(hash []byte) (*token.RefreshToken, error) {
	var token token.RefreshToken
	query := `SELECT id, expires_at, revoked, user_id, family_id FROM tokens WHERE token_hash = $1`
	err := repo.db.Get(&token, query, hash)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (repo *tokenRepository) FindUserIdByHash(hash []byte) (int, error) {
	var token token.RefreshToken
	query := `SELECT user_id FROM tokens WHERE token_hash = $1`
	err := repo.db.Get(&token, query, hash)
	if err != nil {
		return 0, err
	}
	return token.UserId, nil
}

func (repo *tokenRepository) Revoke(hash []byte, uId int) error {
	t := &token.RefreshToken{
		UserId: uId,
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