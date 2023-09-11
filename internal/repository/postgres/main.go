package postgres

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(ctx context.Context) *pgxpool.Pool {
	log.Println("[postgres-pool] init...")

	connStr := os.Getenv("DATABASE_URL")

	if connStr == "" {
		log.Fatalf("[postgres-pool] connection string is empty")
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("[postgres-pool] init error: %s", err)
	}

	log.Println("[postgres-pool] check conn")

	conn, err := pool.Acquire(ctx)
	if err != nil {
		log.Fatalf("[postgres-pool] check conn error: %s", err)
	}

	conn.Release()
	log.Println("[postgres-pool] check conn OK")
	log.Println("[postgres-pool] init done")

	log.Println("[postgres-pool] set time zone: Asia/Almaty")
	if _, err = pool.Exec(ctx, "SET TIME ZONE 'Asia/Almaty'"); err != nil {
		log.Fatalf("[postgres-pool] set time zone error: %s", err)
	}

	return pool
}
