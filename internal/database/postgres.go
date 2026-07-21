package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect(databaseURL string) {

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	// Disable prepared statement cache (required for PgBouncer/Supabase Pooler)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	if err := pool.Ping(context.Background()); err != nil {
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