package postgres

import (
	"context"
	"fmt"
	"session_manager/internal/domain"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	CreateSession(ctx context.Context, sess *domain.Session) error
	PingSession(ctx context.Context, ps *domain.PingSession) error
	GetOnlineSessions(ctx context.Context) (domain.Sessions, error)
}

func NewStorage(pool *pgxpool.Pool) Storage {
	return &storage{pool: pool}
}

type storage struct {
	pool *pgxpool.Pool
}

var (
	createSessionQuery = `INSERT INTO sessions (id, comp_name, ip_addr, login, date_time, createdAt) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	pingSessionQuery = `INSERT INTO sessions_ping (session_id, session_type, date_time, createdAt)
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (session_id)
	DO UPDATE SET
		date_time = EXCLUDED.date_time
		updatedAt = current_timestamp;`

	onlineSessionsQuery = `SELECT s.comp_name, s.ip_addr, s.login, s.date_time AS started_at, sp.date_time AS updated_in
	FROM sessions s
	INNER JOIN sessions_ping sp ON s.id = sp.id
	WHERE`
)

func (s *storage) CreateSession(ctx context.Context, sess *domain.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx, createSessionQuery,
		sess.ID,
		sess.ComputerName,
		sess.IPAddress,
		sess.Login,
		sess.Status,
		sess.DateTime,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("CreateSession: %w", err)
	}
	return nil
}

func (s *storage) PingSession(ctx context.Context, ps *domain.PingSession) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := s.pool.Exec(ctx, pingSessionQuery,
		ps.SessionID,
		ps.DateTime,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("PingSession: %w", err)
	}
	return nil
}

func (s *storage) GetOnlineSessions(ctx context.Context) (domain.Sessions, error) {
	return nil, nil
}
