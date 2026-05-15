package main

import (
	"fmt"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/jmoiron/sqlx"
)

func SeedPosts(db *sqlx.DB, role string, count int) error {
	var userIDs []int
	query := `SELECT u.id
		FROM users AS u
		JOIN roles AS r ON r.id = u.role_id
		WHERE r.role = $1
		`
	err := db.Select(&userIDs, query, role)
	if err != nil {
		return err
	}
	if len(userIDs) == 0 {
		return fmt.Errorf("No users with role %s", role)
	}

	for i := 1; i <= count; i++ {
		title := "What is Lorem Ipsum?"
		excerpt := "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book."
		content := "<h2>What is Lorem Ipsum?</h2><p>Lorem Ipsum is simply dummy text of the printing and typesetting industry. </p><p>Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. </p><h3>What is Lorem Ipsum?</h3><p>It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. </p><p>It was popularised in the <strong>1960s with the release of Letraset sheets</strong> containing Lorem Ipsum passages, and more recently with desktop publishing software like <em>Aldus PageMaker</em> including versions of Lorem Ipsum.</p>"

		userID := userIDs[0]

		var id int
		q := `INSERT INTO posts (title, content, excerpt, user_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`
		err := db.Get(&id, q, title, content, excerpt, userID)
		if err != nil {
			return err
		}

		pSlug := slug.Make(title) + "-" + strconv.Itoa(id)

		_, err = db.Exec(`UPDATE posts SET slug = $1 WHERE id = $2`, pSlug, id)
		if err != nil {
			return err
		}
	}
	fmt.Println("Seed posts created ✔️")
	return nil
}
