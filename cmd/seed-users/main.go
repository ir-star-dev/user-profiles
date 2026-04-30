package main

import (
	"flag"
	"fmt"
	"log"
	"user-profiles/configs"
	"user-profiles/internal/storage/db"

	"github.com/jmoiron/sqlx"
)

func main() {
	// Config
	conf, err := configs.Load()
	if err != nil {
		log.Println("Failed to load config: %w", err)
	}

	// DB
	dbConn, err := db.Connect(conf)
	if err != nil {
		log.Println("Failed to connect db: %w", err)
	}
	defer dbConn.Close()

	role := flag.String("role", "user", "Create test users with the role: default `user`")
	count := flag.Int("count", 30, "Number of test data: default 30")
	flag.Parse()

	err = seedUsers(dbConn, *role, *count)
	if err != nil {
		log.Println("failed to create test seeds %w", err)
	}
}

func seedUsers(db *sqlx.DB, role string, count int) error {
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
