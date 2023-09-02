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
	createSessionQuery = `INSERT INTO sessions (id, comp_name, ip_addr, login, date_time) 
	VALUES ($1, $2, $3, $4, $5);`

	pingSessionQuery = `INSERT INTO sessions_ping (session_id, session_type, date_time)
	VALUES ($1, $2, $3)
	ON CONFLICT (session_id, session_type)
	DO UPDATE SET
		date_time = EXCLUDED.date_time
		updatedAt = current_timestamp;`

	lastSessionQuery = `INSERT INTO sessions_last (id, comp_name, ip_addr, login, date_time)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (comp_name)
	DO UPDATE SET
		session_id = EXCLUDED.session_id,
		ip_addr = EXCLUDED.ip_addr,
		login = EXCLUDED.login,
		date_time = EXCLUDED.date_time;`

	onlineSessionsQuery = `SELECT id, comp_name, ip_addr, login, date_time
	FROM sessions_last;`
)

func (s *storage) CreateSession(ctx context.Context, sess *domain.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := s.pool.Exec(ctx, createSessionQuery,
		sess.ID,
		sess.ComputerName,
		sess.IPAddress,
		sess.Login,
		sess.DateTime,
	); err != nil {
		return err
	}

	if _, err := s.pool.Exec(ctx, lastSessionQuery,
		sess.ID,
		sess.ComputerName,
		sess.IPAddress,
		sess.Login,
		sess.DateTime,
	); err != nil {
		return fmt.Errorf("last: %w", err)
	}

	return nil
}

func (s *storage) PingSession(ctx context.Context, ps *domain.PingSession) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := s.pool.Exec(ctx, pingSessionQuery,
		ps.SessionID,
		ps.SessionType,
		ps.DateTime,
	); err != nil {
		return err
	}

	if _, err := s.pool.Exec(ctx, lastSessionQuery,
		ps.SessionID,
		ps.ComputerName,
		ps.IPAddress,
		ps.Login,
		ps.DateTime,
	); err != nil {
		return fmt.Errorf("last: %w", err)
	}

	return nil
}

func (s *storage) GetOnlineSessions(ctx context.Context) (domain.Sessions, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx, onlineSessionsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make(domain.Sessions, 0, 250)

	for rows.Next() {
		s := domain.Session{}
		if err := rows.Scan(
			&s.ID,
			&s.ComputerName,
			&s.IPAddress,
			&s.Login,
			&s.DateTime,
		); err != nil {
			return nil, fmt.Errorf("iterate row: %w", err)
		}
		sessions = append(sessions, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows all: %w", err)
	}

	return sessions, nil
}
