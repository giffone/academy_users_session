package server

import (
	"context"
	"log"
	"session_manager/internal/repository/postgres"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Env struct {
	pool *pgxpool.Pool
}

func NewEnv(ctx context.Context) *Env {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return &Env{
		pool: postgres.NewPostgres(ctx),
	}
}

func (e *Env) Stop(ctx context.Context) {
	// ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()

	e.pool.Close()
	log.Println("envorinments stopped")
}
