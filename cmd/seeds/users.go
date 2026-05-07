package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

func SeedUsers(db *sqlx.DB, role string, count int) error {
	var prefix string
	if role == "admin" {
		prefix = "admin"
	} else {
		prefix = "test"
	}
	password := "$2a$10$BekeuI1csXF3zQXLRsgecu3W3zG75Bx82P0H5OajMz2ksb5.0IcdG"

	query := `
		INSERT INTO users (name, email, password, role_id)
		VALUES ($1, $2, $3, (SELECT id FROM roles WHERE role = $4))
	`

	for i := 1; i <= count; i++ {
		name := fmt.Sprintf("%s%d", prefix, i)
		email := fmt.Sprintf("%s%d@test.com", prefix, i)

		_, err := db.Exec(query, name, email, password, role)
		if err != nil {
			return err
		}
	}
	log.Printf("Seed users with role %s created ✔️", role)
	return nil
}
