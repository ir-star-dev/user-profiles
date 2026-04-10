package user

type Repository interface {
	Create(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
	FindById(uId int) (*User, error)
	Delete(uId int) error
	UpdateName(user *User) (*User, error)
	FindEmailById(uId int) (string, error)
}

