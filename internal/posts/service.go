package posts

import (
	"errors"
	"time"
	"user-profiles/internal/models"
	"user-profiles/internal/validator"
)

type postService struct {
	pRepo Repository
}

func NewPostService(pRepo Repository) PostService {
	return &postService{pRepo: pRepo}
}

func (s *postService) Create(title, content, excerpt string, uId int) ([]validator.FormValidationErr, error) {
	formValiErr, err := ValidatePostForm(title, content, excerpt)
	if err != nil {
		return formValiErr, errors.New(InvalidForm)
	}

	post := &models.Post{
		Title:   title,
		Content: content,
		Excerpt: excerpt,
		UserId:  uId,
	}
	_, err = s.pRepo.Create(post)
	if err != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "form",
			Message: SomethingWrong,
		})
		return formValiErr, errors.New(SomethingWrong)
	}
	return nil, nil
}

func (s *postService) Update(title, content, excerpt, createdAt string, pId int) ([]validator.FormValidationErr, error) {
	formValiErr, err := ValidatePostForm(title, content, excerpt)
	if err != nil {
		return formValiErr, errors.New(InvalidForm)
	}

	t, err := time.Parse("2006-01-02 15:04:05.999999 -0700 MST", createdAt)
	if err != nil {
		t = time.Now()
	}
	post := &models.UpdatePostRequest{
		Id:        pId,
		Title:     title,
		Excerpt:   excerpt,
		Content:   content,
		Approved:  false,
		CreatedAt: t,
		UpdatedAt: time.Now(),
	}
	_, err = s.pRepo.Update(post)
	if err != nil {
		formValiErr = append(formValiErr, validator.FormValidationErr{
			Name:    "form",
			Message: SomethingWrong,
		})
		return formValiErr, errors.New(SomethingWrong)
	}
	return nil, nil
}

func (s *postService) Delete(pId int) error {
	err := s.pRepo.Delete(pId)
	if err != nil {
		return err
	}
	return nil
}

func (s *postService) FindById(pId int) (*models.PostWithUserName, error) {
	posts, err := s.pRepo.FindById(pId)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *postService) FindByUsername(username string) ([]models.PostWithUserName, error) {
	var posts []models.PostWithUserName
	posts, err := s.pRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *postService) GetOnPage(page int, limit int, approved *bool, username string) ([]models.PostWithUserName, int, error) {
	posts, total, err := s.pRepo.GetOnPage(page, limit, approved, username)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (s *postService) PostsStatus() ([]models.PostsStatus, error) {
	posts, err := s.pRepo.PostsStatus()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *postService) Review(pId int) error {
	err := s.pRepo.Review(pId)
	if err != nil {
		return err
	}
	return nil
}

func (s *postService) Publish(pId int) error {
	err := s.pRepo.Publish(pId)
	if err != nil {
		return err
	}
	return nil
}

func (s *postService) Authors() ([]models.Authors, error) {
	authors, err := s.pRepo.Authors()
	if err != nil {
		return nil, err
	}
	return authors, nil
}

func (s *postService) CreatedPosts(uId int) ([]models.PostRows, error) {
	posts, err := s.pRepo.GetMyPostList(uId)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *postService) PreviewPost(pId int) (*models.PostWithUserName, error) {
	post, err := s.pRepo.FindById(pId)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func ValidatePostForm (title, content, excerpt string) ([]validator.FormValidationErr, error) {
	var formValiErr []validator.FormValidationErr
	var v validator.Validator
	v.CheckField(validator.MaxChars(title, 99), "title", LongTitle)
	v.CheckField(validator.NotBlank(title), "title", EmptyTitle)
	v.CheckField(validator.MinChars(excerpt, 100), "excerpt", ShortExcerpt)
	v.CheckField(validator.MaxChars(excerpt, 300), "excerpt", LongExcerpt)
	v.CheckField(validator.NotBlank(excerpt), "excerpt", EmptyExcerpt)
	v.CheckField(validator.NotBlank(content), "content", EmptyContent)

	if !v.Valid() {
		for key, value := range v.FieldErrors {
			formValiErr = append(formValiErr, validator.FormValidationErr{
				Name:    key,
				Message: value,
			})
		}
		return formValiErr, errors.New(InvalidForm)
	}
	return nil, nil
}