package auth_postgres

import (
	"user-profiles/internal/auth"
	"user-profiles/internal/models"

	"github.com/jmoiron/sqlx"
)

type loginsRepository struct {
	db *sqlx.DB
}

func NewLoginsRepository(db *sqlx.DB) auth.LoginsRepository {
	return &loginsRepository{db: db}
}

func (repo *loginsRepository) Save(login *models.Logins) error {
	query := `
        INSERT INTO logins (user_id, device, ip, user_agent)
		VALUES (:user_id, :device, :ip, :user_agent)
    `
	_, err := repo.db.NamedExec(query, login)
	if err != nil {
		return err
	}
	return nil
}

func (repo *loginsRepository) Get() (*[]models.LoginsResponse, error) {
	var res []models.LoginsResponse
	query := `
        SELECT l.*, u.username  FROM logins AS l
		JOIN users AS u ON l.user_id = u.id
		ORDER BY l.id DESC
		LIMIT 7
    `
	err := repo.db.Select(&res, query)
	if err != nil {
		return nil, err
	}
	return &res, nil
}