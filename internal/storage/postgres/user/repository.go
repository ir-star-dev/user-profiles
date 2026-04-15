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
        INSERT INTO users (name, email, password, role_id, banned)
		SELECT 
			:name,
			:email,
			:password,
			:banned
			r.id
		FROM roles AS r
		WHERE r.role = :role
    `
	_, err := repo.db.NamedExec(query, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) FindById(uId int) (*user.User, error) {
	query := `
		SELECT 
			u.name,
			u.email,
			u.banned,
			r.role
		FROM users AS u
		JOIN roles AS r ON r.id = u.role_id
		WHERE u.id = $1
	`
	var user user.User
	err := repo.db.Get(&user, query, uId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) FindRoleByUserId(uId int) (string, error) {
	query := `
		SELECT role 
		FROM roles as r 
		JOIN users 
		as u ON r.id = u.role_id 
		WHERE u.id = $1
	`
	var role string
	err := repo.db.Get(role, query, uId)
	if err != nil {
		return "", err
	}
	return role, nil
}

func (repo *userRepository) FindByEmail(email string) (*user.User, error) {
	query := `
		SELECT 
			u.id,
			u.banned, 
			u.password,
			r.role
		FROM users AS u
		JOIN roles AS r ON r.id = u.role_id
		WHERE u.email = $1
	`
	var user user.User
	err := repo.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) FindBanStatus(uId int) (*bool, error) {
	query := `
		SELECT banned 
		FROM users 
		WHERE id = $1
	`
	var status bool
	err := repo.db.Get(status, query, uId)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (repo *userRepository) Delete(uId int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := repo.db.Exec(query, uId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) UpdateName(userName string, uId int) (*user.User, error) {
	user := &user.User{
		Id:   uId,
		Name: userName,
	}
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

func (repo *userRepository) GetAll() ([]*user.User, error) {
	query := `
		SELECT 
			u.id, 
			u.name,
			u.banned, 
			u.email, 
			u.created_at, 
			r.role
		FROM users AS u
		JOIN roles AS r ON r.id = u.role_id
	`
	users := []*user.User{}
	err := repo.db.Select(&users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) Ban(uId int) error {
	user := &user.User{
		Id:     uId,
		Banned: "true",
	}
	query := `
        UPDATE users 
        SET banned = :banned
        WHERE id = :id
    `
	_, err := repo.db.NamedExec(query, user)
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) Unban(uId int) error {
	user := &user.User{
		Id:     uId,
		Banned: "false",
	}
	query := `
        UPDATE users 
        SET banned = :banned
        WHERE id = :id
    `
	_, err := repo.db.NamedExec(query, user)
	if err != nil {
		return err
	}
	return nil
}
