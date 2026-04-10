package main

import (
	"user-profiles/configs"
	"log"
	"flag"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	up := flag.Bool("up", false, "Apply all migrations")
	down := flag.Bool("down", false, "Rollback all migrations")
	flag.Parse()

	if !*up && !*down {
		log.Println("Usage:")
		log.Println("  -up     Apply migrations")
		log.Println("  -down   Rollback all migrations")
		os.Exit(1)
	}

    conf, err := configs.Load()
	if err != nil {
        log.Fatal(err)
	}
	dsn := conf.Db.Dsn
    mPath := conf.Db.MigrationsPath

	
	if *up {
		if err := MigrationsUp(dsn, mPath); err != nil {
			log.Fatal(err)
		}
	}

	if *down {
		if err := MigrationsDownAll(dsn, mPath); err != nil {
			log.Fatal(err)
		}
	}
}

func MigrationsUp(dsn string, mPath string) error {
	m, err := migrate.New(
		"file://"+mPath,
		dsn,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations applied successfully ✔️")
	return nil
}

func MigrationsDownAll(dsn string, mPath string) error {
    m, err := migrate.New(
        "file://"+mPath,
        dsn,
    )
    if err != nil {
        return err
    }

    if err := m.Down(); err != nil && err != migrate.ErrNoChange {
        return err
    }

    log.Println("All migrations rolled back ✔️")
	return nil
}