package profile_service

import (
	"user-profiles/internal/domain/user"
)

type Service interface {
	View(uId int) (*user.UserWithRole, error)
	ViewAll() ([]*user.UserWithRole, error)
	ChangeName(uId int, name *string) (*user.UserWithRole, error)
	Delete(uId int) error
	Ban(uId int) error
	Unban(uId int) error
}