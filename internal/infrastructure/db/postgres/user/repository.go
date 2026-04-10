package postgres

import (
	"user-profiles/internal/domain/user"

	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) user.Repository {
	return &userRepository{db: db}
}

func (repo *userRepository) Create(user *user.User) (*user.User, error) {
	query := `
        INSERT INTO users (name, email, password) 
        VALUES (:name, :email, :password)
    `
	_, err := repo.db.NamedExec(query, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) FindById(uId int) (*user.User, error) {
	query := `SELECT * FROM users WHERE id = $1`
	var user user.User
	err := repo.db.Get(&user, query, uId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) FindByEmail(email string) (*user.User, error) {
	query := `SELECT * FROM users WHERE email = $1`
	var user user.User
	err := repo.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) Delete(uId int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := repo.db.Exec(query, uId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) UpdateName(user *user.User) (*user.User, error) {
	query := `
        UPDATE users 
        SET name = :name
        WHERE id = :id
    `
	_, err := repo.db.NamedExec(query, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) FindEmailById(uId int) (string, error) {
	query := `SELECT email FROM users WHERE id = $1`
	var user user.User
	err := repo.db.Get(&user, query, uId)
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
