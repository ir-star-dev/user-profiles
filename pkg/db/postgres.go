package db

import (
	"user-profiles/configs"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(conf *configs.Config) (*sqlx.DB, error) {
	dsn := conf.Db.Dsn
	db, err := sqlx.Connect(conf.Db.Driver, dsn)
	if err != nil {
		return nil, err
	}
	log.Println("Connected to PostgreSQL ✔️")
	return db, nil
}