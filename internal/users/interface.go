package users

type Repository interface {
	Create(user *UserWithRole) (*UserWithRole, error)
	Delete(uId int) error
	UpdateName(userName string, uId int) (*UserWithRole, error)

	FindByEmail(email string) (*UserWithRole, error)
	FindById(uId int) (*UserWithRole, error)
	FindRoleByUserId(uId int) (string, error)
	FindBanStatus(uId int) (*bool, error)

	GetAll() ([]*UserWithRole, error)
	Ban(uId int) error
	Unban(uId int) error
}

type UsersService interface {
	View(uId int) (*UserWithRole, error)
	ChangeName(uId int, name *string) (*UserWithRole, error)
	Delete(uId int) error
	ViewAll() ([]*UserWithRole, error)
}

type AdminService interface {
	//ViewAll() ([]*UserWithRole, error)
	ChangeNameById(uId int, name *string) (*UserWithRole, error)
	DeleteById(uId int) error
	BanById(uId int) error
	UnbanById(uId int) error
}

