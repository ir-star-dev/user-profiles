package posts

type Repository interface {
	Create(post *Post) (*Post, error)
	Delete(pId int) error
	Update(post *Post) (int64, error)

	FindById(uId int) ([]Post, error)

	GetAll(page int) ([]Post, int, error)
}

type PostService interface {
	Create(post *Post) (*Post, error)
	Delete(pId int) error
	Update(post *Post) (int64, error)

	FindById(uId int) ([]Post, error)

	GetAll(page int) ([]Post, int, error)
}