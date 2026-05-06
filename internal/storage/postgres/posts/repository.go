package posts_postgres

import (
	"user-profiles/cmd/user-profiles/posts"

	"github.com/jmoiron/sqlx"
)

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) posts.Repository {
	return &postRepository{db: db}
}

func (repo *postRepository) Create(post *posts.Post) (*posts.Post, error) {
	query := `
        INSERT INTO posts (title, content, created_at, updated_at, approved, user_id)
		SELECT
			:title,
			:content,
			:created_at,
			:updated_at,
			:approved,
			u.id
		FROM users AS u
		WHERE u.id = :user_id
    `
	_, err := repo.db.NamedExec(query, post)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *postRepository) Delete(pId int) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := repo.db.Exec(query, pId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *postRepository) Update(post *posts.Post) (int64, error) {
	query := `
        UPDATE posts 
        SET title = :title, 
			content = :content, 
			created_at = :created_at, 
			updated_at = :updated_at, 
			approved = :approved,
        WHERE id = :id
    `
	r, err := repo.db.NamedExec(query, post)
	if err != nil {
		return 0, err
	}
	id, err := r.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (repo *postRepository) FindById(uId int) ([]posts.Post, error) {
	var posts []posts.Post
	query := `
		SELECT * FROM posts 
		WHERE user_id = $1
	`
	err := repo.db.Get(&posts, query, uId)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *postRepository) GetAll(page int) ([]posts.Post, int, error) {
	limit := 10
	offset := (page - 1) * limit

	var totalCount int
	countQuery := `
		SELECT COUNT(*)
		FROM posts
	`
	err := repo.db.Get(&totalCount, countQuery)
	if err != nil {
		return nil, 0, err
	}
	query := `
		SELECT *
		FROM posts
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	posts := []posts.Post{}
	err = repo.db.Select(&posts, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return posts, totalCount, nil
}
