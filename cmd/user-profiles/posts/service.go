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

func (s *postService) GetAll(page int, limit int) ([]PostWithUserName, int, error) {
	posts, total, err := s.pRepo.GetAll(page, limit)
	if err != nil {
		return nil, 0, err
	}
	return posts, total, nil
}