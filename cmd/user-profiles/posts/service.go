package posts

type postService struct {
	pRepo Repository
}

func NewPostService(pRepo Repository) PostService {
	return &postService{pRepo: pRepo}
}

func (s *postService) Create(post *Post) (*Post, error) {
	return post, nil
}

func (s *postService) Delete(pId int) error {
	err := s.pRepo.Delete(pId)
	if err != nil {
		return err
	}
	return nil
}

func (s *postService) Update(post *Post) (int64, error) {
	return 0, nil
}

func (s *postService) FindById(pId int) (*PostWithUserName, error) {
	posts, err := s.pRepo.FindById(pId)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *postService) FindByUsername(username string) ([]PostWithUserName, error) {
	var posts []PostWithUserName
	posts, err := s.pRepo.FindByUsername(username)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *postService) GetOnPage(page int, limit int, approved *bool, username string) ([]PostWithUserName, int, error) {
	posts, total, err := s.pRepo.GetOnPage(page, limit, approved, username)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}

func (s *postService) GetAll() ([]Post, error) {
	posts, err := s.pRepo.GetAll()
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (s *postService) PostsStatus() ([]PostsStatus, error) {
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