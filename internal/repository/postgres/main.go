package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func newStorage() *pgxpool.Pool {
	// return nil
	ctx := context.Background()
	log.Println("[postgres-pool] init...")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		// conf.UserName,
		// conf.Password,
		// conf.Host,
		// conf.Port,
		// conf.Name,
	)

	// connStr = "user=postgres password=postgres port=5432 dbname=postgres sslmode=disable"
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