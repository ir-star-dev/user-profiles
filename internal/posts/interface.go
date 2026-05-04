package posts

type Repository interface {
	Create(post *Post) (*Post, error)
	Delete(pId int) error
	Update(post *Post) (int, error)

	FindById(uId int) ([]Post, error)

	GetAll() ([]Post, error)
}

type PostService interface {
	Create(post *Post) (*Post, error)
	Delete(pId int) error
	Update(post *Post) (int, error)

	FindById(uId int) ([]Post, error)

	GetAll() ([]Post, error)
}