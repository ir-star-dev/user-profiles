package user

type Repository interface {
	Create(user *User) (*User, error)
	Delete(uId int) error
	UpdateName(userName string, uId int) (*User, error)

	FindByEmail(email string) (*User, error)
	FindById(uId int) (*User, error)	
	FindEmailById(uId int) (string, error)
	FindRoleByUserId(uId int) (string, error)

	GetAll() ([]*User, error)
}

