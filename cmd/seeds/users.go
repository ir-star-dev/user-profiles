package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func SeedUsers(db *sqlx.DB, role string, count int) error {
	var prefix string
	switch role {
	case "admin": 
		prefix = "admin"
	case "moderator":
		prefix = "mod"
	case "user":
		prefix = "user"
	}
	//admin
	password := "$2a$10$BekeuI1csXF3zQXLRsgecu3W3zG75Bx82P0H5OajMz2ksb5.0IcdG"

	query := `
		INSERT INTO users (name, email, username, password, role_id)
		VALUES ($1, $2, $3, $4, (SELECT id FROM roles WHERE role = $5))
	`

	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("%s%d", prefix, i)
		email := fmt.Sprintf("%s%d@example.com", prefix, i)
		username := fmt.Sprintf("%s%d", prefix, i)

		_, err := db.Exec(query, name, email, username, password, role)
		if err != nil {
			return err
		}
	}
	fmt.Printf("Seed users with role %s created ✔️", role)
	return nil
}
