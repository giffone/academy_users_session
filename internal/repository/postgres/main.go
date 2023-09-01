package postgres

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres() *pgxpool.Pool {
	ctx := context.Background()
	log.Println("[postgres-pool] init...")

	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		log.Fatalf("[postgres-pool] connection string is empty")
	}
	
	pg, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("[postgres-pool] init error: %s", err)
	}

	log.Println("[postgres-pool] check conn")

	conn, err := pg.Acquire(ctx)
	if err != nil {
		log.Fatalf("[postgres-pool] check conn error: %s", err)
	}

	conn.Release()
	log.Println("[postgres-pool] check conn OK")
	log.Println("[postgres-pool] init done")

	return pg
}
