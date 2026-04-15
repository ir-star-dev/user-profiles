package user

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

