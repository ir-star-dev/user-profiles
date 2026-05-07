package main

import (
	"flag"
	"log"
	"strings"
	"user-profiles/configs"
	"user-profiles/internal/storage/db"
)

func main() {
	// Config
	conf, err := configs.Load()
	if err != nil {
		log.Fatal("Failed to load config: %w", err)
	}

	// DB
	dbConn, err := db.Connect(conf)
	if err != nil {
		log.Fatal("Failed to connect db: %w", err)
	}
	defer dbConn.Close()

	mode := flag.String("mode", "", "roles | users | posts")
	roles := flag.String("roles", "admin,moderator,user", "admin,moderator,user")
	role := flag.String("role", "user", "role of user")
	count := flag.Int("count", 30, "number of users or posts")
	flag.Parse()

	switch *mode {

	case "roles":
		if *roles == "" {
			log.Fatal("roles required")
		}
		roleList := strings.Split(*roles, ",")
		err = SeedRoles(dbConn, roleList)
		if err != nil {
			log.Fatal(err)
		}

	case "users":
		err = EnsureRoleExists(dbConn, *role)
		if err != nil {
			log.Fatal(err)
		}

		err = SeedUsers(dbConn, *role, *count)
		if err != nil {
			log.Fatal(err)
		}

	case "posts":
		err = EnsureRoleExists(dbConn, *role)
		if err != nil {
			log.Fatal(err)
		}

		err = SeedPosts(dbConn, *role, *count)
		if err != nil {
			log.Fatal(err)
		}

	default:
		log.Fatal("use -mode=roles or -mode=users or -mode=posts")
	}

	if err != nil {
		log.Println("failed to create test seeds %w", err)
	}
}
