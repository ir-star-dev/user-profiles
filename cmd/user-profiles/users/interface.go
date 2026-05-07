package users

type Repository interface {
	Create(user *UserWithRole) (*UserWithRole, error)
	Delete(uId int) error
	UpdateName(userName string, uId int) (*UserWithRole, error)

	FindByEmail(email string) (*UserWithRole, error)
	FindById(uId int) (*UserWithRole, error)
	FindRoleByUserId(uId int) (string, error)
	FindBanStatus(uId int) (*bool, error)

	GetAll(page int) ([]UserWithRole, int, error)
	Ban(uId int) error
	Unban(uId int) error
}

type UsersService interface {
	View(uId int) (*UserWithRole, error)
	ChangeName(uId int, name *string) (*UserWithRole, error)
	Delete(uId int) error
	ViewAll(page int) ([]UserWithRole, int, error)
	Role(uId int) (string, error)

	Ban(uId int) error
	Unban(uId int) error
}
