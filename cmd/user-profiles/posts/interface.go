package posts

type Repository interface {
	Create(post *Post) (*Post, error)
	Delete(pId int) error
	Update(post *Post) (int64, error)

	FindById(pId int) (*PostWithUserName, error)
	FindByUsername(username string) ([]PostWithUserName, error)

	GetOnPage(page int, limit int) ([]PostWithUserName, int, error)
	GetAll() ([]Post, error)
}

type PostService interface {
	Create(post *Post) (*Post, error)
	Delete(pId int) error
	Update(post *Post) (int64, error)

	FindById(pId int) (*PostWithUserName, error)
	FindByUsername(username string) ([]PostWithUserName, error)

	GetOnPage(page int, limit int) ([]PostWithUserName, int, error)
	GetAll() ([]Post, error)
}