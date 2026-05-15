package posts

type Repository interface {
	Create(post *Post) (*Post, error)
	Delete(pId int) error
	Update(post *UpdatePostRequest) (int64, error)

	Review(pId int) error
	Publish(pId int) error

	FindById(pId int) (*PostWithUserName, error)
	FindByUsername(username string) ([]PostWithUserName, error)

	GetOnPage(page int, limit int, approved *bool, username string) ([]PostWithUserName, int, error)
	GetAll() ([]Post, error)
	PostsStatus() ([]PostsStatus, error)
	
	Authors() ([]Authors, error)
}

type PostService interface {
	//Create(post *Post) (*Post, error)
	Delete(pId int) error
	//Update(post *UpdatePostRequest) (int64, error)

	Review(pId int) error
	Publish(pId int) error

	FindById(pId int) (*PostWithUserName, error)
	FindByUsername(username string) ([]PostWithUserName, error)

	GetOnPage(page int, limit int, approved *bool, username string) ([]PostWithUserName, int, error)
	GetAll() ([]Post, error)
	PostsStatus() ([]PostsStatus, error)
}