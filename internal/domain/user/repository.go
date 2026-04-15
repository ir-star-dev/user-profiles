package user

type Repository interface {
	Create(user *User) (*User, error)
	Delete(uId int) error
	UpdateName(userName string, uId int) (*User, error)

	FindByEmail(email string) (*User, error)
	FindById(uId int) (*User, error)	
	FindRoleByUserId(uId int) (string, error)
	FindBanStatus(uId int) (*bool, error)

	GetAll() ([]*User, error)
	Ban(uId int) error
	Unban(uId int) error
}

