package posts

import (
	"user-profiles/internal/models"
	"user-profiles/internal/validator"
)

type Repository interface {
	Create(post *models.Post) (*models.Post, error)
	Delete(pId int) error
	Update(post *models.UpdatePostRequest) (int64, error)

	Review(pId int) error
	Publish(pId int) error

	FindById(pId int) (*models.PostWithUserName, error)
	FindByUsername(username string) ([]models.PostWithUserName, error)

	GetOnPage(page int, limit int, approved *bool, username string, search string) ([]models.PostWithUserName, int, error)
	PostsStatus() ([]models.PostsStatus, error)

	Authors() ([]models.Authors, error)
	GetMyPostList(uId int) ([]models.PostRows, error)
}

type PostService interface {
	Create(title, content, excerpt string, uId int) ([]validator.FormValidationErr, error)
	Delete(pId int) error
	Update(title, content, excerpt, createdAt string, pId int) ([]validator.FormValidationErr, error)

	Review(pId int) error
	Publish(pId int) error

	FindById(pId int) (*models.PostWithUserName, error)
	FindByUsername(username string) ([]models.PostWithUserName, error)

	GetOnPage(page int, limit int, approved *bool, username string, search string) ([]models.PostWithUserName, int, error)
	PostsStatus() ([]models.PostsStatus, error)

	Authors() ([]models.Authors, error)
	CreatedPosts(uId int) ([]models.PostRows, error)
	PreviewPost(pId int) (*models.PostWithUserName, error)
}
