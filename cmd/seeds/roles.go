package main

import (
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

func SeedRoles(db *sqlx.DB, roles []string) error {
	query := `
		INSERT INTO roles (role)
		VALUES ($1)
		ON CONFLICT (role) DO NOTHING
	`

	for _, role := range roles {
		role = strings.TrimSpace(role)

		if role == "" {
			continue
		}

		_, err := db.Exec(query, role)
		if err != nil {
			return err
		}
	}

	log.Printf("Roles created successfully ✔️")
	return nil
}

func EnsureRoleExists(db *sqlx.DB, role string) error {
	query := `
		INSERT INTO roles (role)
		VALUES ($1)
		ON CONFLICT (role) DO NOTHING
	`

	_, err := db.Exec(query, role)
	if err != nil {
		return err
	}

	log.Printf("Role %s checked ✔️", role)
	return nil
}
