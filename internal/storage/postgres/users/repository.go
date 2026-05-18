package users_postgres

import (
	"strconv"
	"strings"
	"user-profiles/internal/models"
	"user-profiles/internal/users"

	"github.com/jmoiron/sqlx"
)

type usersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) users.Repository {
	return &usersRepository{db: db}
}

func (repo *usersRepository) Create(user *models.UserWithRole) (*models.UserWithRole, error) {
	query := `
        INSERT INTO users (name, email, username, password, role_id)
		SELECT 
			:name,
			:email,
			:username,
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

func (repo *usersRepository) FindById(uId int) (*models.UserWithRole, error) {
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
	var user models.UserWithRole
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

func (repo *usersRepository) FindByEmail(email string) (*models.UserWithRole, error) {
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
	var user models.UserWithRole
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

func (repo *usersRepository) UpdateName(userName string, uId int) (*int, error) {
	user := &models.UserWithRole{
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
	return &uId, nil
}

func (repo *usersRepository) GetOnPage(page, limit int, banned *bool, role, search string) ([]models.UserWithRole, int, error) {
	offset := (page - 1) * limit

	args := []any{}
	conditions := []string{}

	// search
	if search != "" {
		args = append(args, "%"+search+"%")
		conditions = append(conditions, "(u.name ILIKE $"+strconv.Itoa(len(args))+" OR u.username ILIKE $"+strconv.Itoa(len(args))+" )")
	}

	// approved
	if banned != nil {
		args = append(args, *banned)
		conditions = append(conditions, "u.banned = $"+strconv.Itoa(len(args)))
	}

	// username
	if role != "" {
		args = append(args, role)
		conditions = append(conditions, "r.role = $"+strconv.Itoa(len(args)))
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// COUNT
	var totalCount int
	countQuery := `
		SELECT COUNT(*)
		FROM users AS u
		JOIN roles AS r ON u.role_id = r.id
	` + where

	err := repo.db.Get(&totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// pagination
	args = append(args, limit, offset)

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
		` + where + `
		ORDER BY u.id DESC
		LIMIT $` + strconv.Itoa(len(args)-1) + `
		OFFSET $` + strconv.Itoa(len(args))

	users := []models.UserWithRole{}
	err = repo.db.Select(&users, query, args...)
	if err != nil {
		return nil, 0, err
	}
	return users, totalCount, nil
}

func (repo *usersRepository) GetAll() ([]models.UserWithRole, error) {
	query := `SELECT * FROM users`
	var users []models.UserWithRole
	err := repo.db.Select(&users, query)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *usersRepository) Ban(uId int) error {
	user := &models.User{
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
	user := &models.User{
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

func (repo *usersRepository) CountRoles() ([]models.StatUsers, error) {
	query := `SELECT 
			r.role,
			COUNT(u.id) AS count,
			COUNT(u.id) FILTER (WHERE u.banned = true) AS banned
		FROM users AS u
		JOIN roles AS r 
			ON u.role_id = r.id
		GROUP BY r.role`
	var stats []models.StatUsers
	err := repo.db.Select(&stats, query)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func (repo *usersRepository) Roles() ([]models.Roles, error) {
	var roles []models.Roles
	query := `SELECT * FROM roles ORDER BY role`
	err := repo.db.Select(&roles, query)
	if err != nil {
		return nil, err
	}
	return roles, nil
}
