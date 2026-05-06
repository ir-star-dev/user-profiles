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

func (s *postService) FindById(uId int) ([]Post, error) {
	var posts []Post
	return posts, nil
}

func (s *postService) GetAll(page int) ([]Post, int, error) {
	var posts []Post
	return posts, 0, nil
}