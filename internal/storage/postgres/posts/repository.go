package posts_postgres

import (
	"strconv"
	"strings"
	"user-profiles/cmd/user-profiles/posts"

	"github.com/jmoiron/sqlx"
	"github.com/gosimple/slug"
)

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) posts.Repository {
	return &postRepository{db: db}
}

func (repo *postRepository) Create(post *posts.Post) (*posts.Post, error) {
	var id int
	query := `INSERT INTO posts (title, content, excerpt, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	err := repo.db.Get(&id, query, post.Title, post.Content, post.Excerpt, post.UserId)
	if err != nil {
		return nil, err
	}

	pSlug := slug.Make(post.Title) + "-" + strconv.Itoa(id)

	_, err = repo.db.Exec(`UPDATE posts SET slug = $1 WHERE id = $2`, pSlug, id)
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

func (repo *postRepository) Update(post *posts.UpdatePostRequest) (int64, error) {
	query := `
        UPDATE posts 
        SET title = :title, 
			content = :content,
			excerpt = :excerpt, 
			created_at = :created_at, 
			updated_at = :updated_at, 
			approved = :approved
        WHERE id = :id
    `
	_, err := repo.db.NamedExec(query, post)
	if err != nil {
		return 0, err
	}
	return int64(post.Id), nil
}

func (repo *postRepository) FindById(pId int) (*posts.PostWithUserName, error) {
	var posts posts.PostWithUserName
	query := `
		SELECT 
			p.*,
			u.username
		FROM posts AS p
		JOIN users AS u 
		ON p.user_id = u.id 
		WHERE p.id = $1
	`
	err := repo.db.Get(&posts, query, pId)
	if err != nil {
		return nil, err
	}
	return &posts, nil
}

func (repo *postRepository) FindByUsername(username string) ([]posts.PostWithUserName, error) {
	query := `
		SELECT 
			p.*,
			u.username
		FROM posts AS p
		JOIN users AS u 
		ON p.user_id = u.id
		WHERE u.username = $1 AND
		p.approved = true
		ORDER BY p.created_at DESC
	`
	p := []posts.PostWithUserName{}
	err := repo.db.Select(&p, query, username)
	if err != nil {
		return nil, err
	}
	if len(p) == 0 {
		return nil, posts.PostNotFound
	}
	return p, nil
}

func (repo *postRepository) GetOnPage(page int, limit int, approved *bool, username string) ([]posts.PostWithUserName, int, error) {
	offset := (page - 1) * limit

	args := []any{}
	conditions := []string{}

	// approved
	if approved != nil {
		args = append(args, *approved)
		conditions = append(conditions, "p.approved = $"+strconv.Itoa(len(args)))
	}

	// username
	if username != "" {
		args = append(args, username)
		conditions = append(conditions, "u.username = $"+strconv.Itoa(len(args)))
	}

	where := ""
	if len(conditions) > 0 {
		where = "WHERE " + strings.Join(conditions, " AND ")
	}

	// COUNT
	var totalCount int

	countQuery := `
		SELECT COUNT(*)
		FROM posts AS p
		JOIN users AS u ON p.user_id = u.id
	` + where

	err := repo.db.Get(&totalCount, countQuery, args...)
	if err != nil {
		return nil, 0, err
	}

	// pagination
	args = append(args, limit, offset)

	query := `
		SELECT 
			p.*,
			u.username
		FROM posts AS p
		JOIN users AS u 
			ON p.user_id = u.id
		` + where + `
		ORDER BY p.created_at DESC
		LIMIT $` + strconv.Itoa(len(args)-1) + `
		OFFSET $` + strconv.Itoa(len(args))

	posts := []posts.PostWithUserName{}

	err = repo.db.Select(&posts, query, args...)
	if err != nil {
		return nil, 0, err
	}

	return posts, totalCount, nil
}

func (repo *postRepository) GetAll() ([]posts.Post, error) {
	query := `SELECT * FROM posts`
	var posts []posts.Post
	err := repo.db.Select(&posts, query)
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (repo *postRepository) PostsStatus() ([]posts.PostsStatus, error) {
	query := `SELECT
			COUNT(id) FILTER (WHERE approved = false) AS pending,
			COUNT(id) FILTER (WHERE approved = true) AS published
		FROM posts
	`
	var pStats []posts.PostsStatus
	err := repo.db.Select(&pStats, query)
	if err != nil {
		return nil, err
	}
	return pStats, nil
}

func (repo *postRepository) Authors() ([]posts.Authors, error) {
	query := `SELECT u.username
		FROM posts AS p
		JOIN users AS u ON p.user_id = u.id
		GROUP BY u.username
	`
	var authors []posts.Authors
	err := repo.db.Select(&authors, query)
	if err != nil {
		return nil, err
	}
	return authors, nil
}

func (repo *postRepository) Review(uId int) error {
	post := &posts.Post{
		Id:       uId,
		Approved: false,
	}
	query := `
        UPDATE posts 
        SET approved = :approved
        WHERE id = :id
    `
	_, err := repo.db.NamedExec(query, post)
	if err != nil {
		return err
	}
	return nil
}

func (repo *postRepository) Publish(uId int) error {
	post := &posts.Post{
		Id:       uId,
		Approved: true,
	}
	query := `
        UPDATE posts 
        SET approved = :approved
        WHERE id = :id
    `
	_, err := repo.db.NamedExec(query, post)
	if err != nil {
		return err
	}
	return nil
}
