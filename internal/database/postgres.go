package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(databaseURL string) {

	pool, err := pgxpool.New(context.Background(), databaseURL)

	if err != nil {
		log.Fatal(err)
	}

	err = pool.Ping(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	DB = pool

	log.Println("Connected to PostgreSQL")
}

func GetDB() *pgxpool.Pool {
	return DB
}

func Close() {

	if DB != nil {
		DB.Close()
	}
}