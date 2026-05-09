package users_postgres

import (
	"user-profiles/cmd/user-profiles/users"

	"github.com/jmoiron/sqlx"
)

type usersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) users.Repository {
	return &usersRepository{db: db}
}

func (repo *usersRepository) Create(user *users.UserWithRole) (*users.UserWithRole, error) {
	query := `
        INSERT INTO users (name, email, password, role_id)
		SELECT 
			:name,
			:email,
			:password,
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

func (repo *usersRepository) FindById(uId int) (*users.UserWithRole, error) {
	query := `
		SELECT 
			u.id,
			u.name,
			u.email,
			u.banned,
			u.created_at,
			r.role
		FROM users AS u
		JOIN roles AS r ON r.id = u.role_id
		WHERE u.id = $1
	`
	var user users.UserWithRole
	err := repo.db.Get(&user, query, uId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *usersRepository) FindRoleByUserId(uId int) (string, error) {
	query := `
		SELECT role 
		FROM roles as r 
		JOIN users 
		as u ON r.id = u.role_id 
		WHERE u.id = $1
	`
	var role string
	err := repo.db.Get(&role, query, uId)
	if err != nil {
		return "", err
	}
	return role, nil
}

func (repo *usersRepository) FindByEmail(email string) (*users.UserWithRole, error) {
	query := `
		SELECT 
			u.id,
			u.password,
			u.banned,
			r.role
		FROM users AS u
		JOIN roles AS r ON r.id = u.role_id
		WHERE u.email = $1
	`
	var user users.UserWithRole
	err := repo.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *usersRepository) FindBanStatus(uId int) (*bool, error) {
	query := `
		SELECT banned 
		FROM users 
		WHERE id = $1
	`
	var status bool
	err := repo.db.Get(&status, query, uId)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (repo *usersRepository) Delete(uId int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := repo.db.Exec(query, uId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *usersRepository) UpdateName(userName string, uId int) (*users.UserWithRole, error) {
	user := &users.UserWithRole{
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

func (repo *usersRepository) GetOnPage(page int) ([]users.UserWithRole, int, error) {
	limit := 10
	offset := (page-1)*limit

	var totalCount int
	countQuery := `
		SELECT COUNT(*)
		FROM users
	`
	err := repo.db.Get(&totalCount, countQuery)
	if err != nil {
		return nil, 0, err
	}

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
		ORDER BY u.id DESC
		LIMIT $1 OFFSET $2
	`
	users := []users.UserWithRole{}
	err = repo.db.Select(&users, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return users, totalCount, nil
}

func (repo *usersRepository) GetAll() ([]users.UserWithRole, error) {
	query := `SELECT * FROM users`
	var users []users.UserWithRole
	err := repo.db.Select(&users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *usersRepository) Ban(uId int) error {
	user := &users.User{
		Id:     uId,
		Banned: true,
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

func (repo *usersRepository) Unban(uId int) error {
	user := &users.User{
		Id:     uId,
		Banned: false,
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

func (repo *usersRepository) CountRoles() ([]users.StatUsers, error) {
	query := `SELECT 
			r.role,
			COUNT(u.id) AS count,
			COUNT(u.id) FILTER (WHERE u.banned = true) AS banned
		FROM users AS u
		JOIN roles AS r 
			ON u.role_id = r.id
		GROUP BY r.role`
	var stats []users.StatUsers
	err := repo.db.Select(&stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}