package users

import (
	"user-profiles/internal/models"
	"user-profiles/internal/validator"
)

type Repository interface {
	Create(user *models.UserWithRole) (*models.UserWithRole, error)
	Delete(uId int) error
	UpdateName(userName string, uId int) (*int, error)

	FindByEmail(email string) (*models.UserWithRole, error)
	FindById(uId int) (*models.UserWithRole, error)
	FindRoleByUserId(uId int) (string, error)
	FindBanStatus(uId int) (*bool, error)

	GetOnPage(page int, limit int, banned *bool, role string) ([]models.UserWithRole, int, error)
	CountRoles() ([]models.StatUsers, error)
	Roles() ([]models.Roles, error)

	Ban(uId int) error
	Unban(uId int) error
}

type UsersService interface {
	CreateUser(email, name, password, role string) ([]validator.FormValidationErr, error)
	Get(uId int) (*models.UserWithRole, error)
	ChangeName(uId int, name string) ([]validator.FormValidationErr, *string, error)
	Delete(uId int) error
	GetOnPage(page int, limit int, banned *bool, role string) ([]models.UserWithRole, int, error)
	Role(uId int) (string, error)

	Ban(uId int) error
	Unban(uId int) error

	CountRoles() ([]models.StatUsers, error)
	GetRoles() ([]models.Roles, error)
}
