package postgres

import (
	"context"
	"fmt"
	"session_manager/internal/domain"
	"session_manager/internal/domain/request"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	CreateSession(ctx context.Context, sess *request.Session) error
	PingSession(ctx context.Context, ps *request.PingSession) error
	GetOnlineDashboard(ctx context.Context) ([]domain.Session, error)
	IsSessionExists(ctx context.Context, comp_name, login string) (*domain.Session, error)
}

func NewStorage(pool *pgxpool.Pool) Storage {
	return &storage{pool: pool}
}

type storage struct {
	pool *pgxpool.Pool
}

const (
	zeroStart  = "zero start"
	zeroEnd    = "zero end"

	createSessionQuery = `INSERT INTO sessions (id, comp_name, ip_addr, login, next_ping_sec, date_time) 
	VALUES ($1, $2, $3, $4, $5, $6);`

	pingSessionQuery = `INSERT INTO sessions_ping (session_id, session_type, date_time)
	VALUES ($1, $2, $3)
	ON CONFLICT (session_id, session_type)
	DO UPDATE SET
		date_time = EXCLUDED.date_time;`

	createOnlineSessionQuery = `INSERT INTO online_dashboard (session_id, comp_name, ip_addr, login, date_time)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (comp_name)
	DO UPDATE SET
		session_id = EXCLUDED.session_id,
		ip_addr = EXCLUDED.ip_addr,
		login = EXCLUDED.login,
		date_time = EXCLUDED.date_time;`

	onlineDashboardQuery = `SELECT session_id, comp_name, ip_addr, login, date_time
	FROM online_dashboard;`

	isOnlineSessionQuery = `SELECT session_id, comp_name, ip_addr, login, date_time
	FROM online_dashboard
	WHERE comp_name = $1
	AND login = $2
	AND date_time >= NOW();`
)

func (s *storage) CreateSession(ctx context.Context, req *request.Session) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	// start session
	if _, err := tx.Exec(ctx, createSessionQuery,
		req.ID,
		req.ComputerName,
		req.IPAddress,
		req.Login,
		req.NextPingSeconds,
		req.DateTime,
	); err != nil {
		return fmt.Errorf("start session: %w", err)
	}

	// put in online dashboard
	if _, err := tx.Exec(ctx, createOnlineSessionQuery,
		req.ID,
		req.ComputerName,
		req.IPAddress,
		req.Login,
		req.DateTime.Add(time.Duration(req.NextPingSeconds)*time.Second),
	); err != nil {
		return fmt.Errorf("online dashboard: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (s *storage) PingSession(ctx context.Context, req *request.PingSession) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	endDate := req.DateTime.Add(time.Duration(req.NextPingSeconds) * time.Second)

	if _, err := s.pool.Exec(ctx, pingSessionQuery,
		req.SessionID,
		req.SessionType,
		endDate,
	); err != nil {
		return err
	}

	// put in online dashboard
	if _, err := s.pool.Exec(ctx, createOnlineSessionQuery,
		req.SessionID,
		req.ComputerName,
		req.IPAddress,
		req.Login,
		endDate,
	); err != nil {
		return fmt.Errorf("online dashboard: %w", err)
	}

	return nil
}

func (s *storage) GetOnlineDashboard(ctx context.Context) ([]domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx, onlineDashboardQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]domain.Session, 0, 250)

	for rows.Next() {
		s := domain.Session{}
		if err := rows.Scan(
			&s.SessionID,
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

func (s *storage) IsSessionExists(ctx context.Context, comp_name, login string) (*domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var session domain.Session

	if err := s.pool.QueryRow(ctx, isOnlineSessionQuery,
		comp_name,
		login,
	).Scan(&session); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}
