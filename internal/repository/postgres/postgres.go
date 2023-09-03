package postgres

import (
	"context"
	"fmt"
	"session_manager/internal/domain"
	"session_manager/internal/domain/request"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	CreateSession(ctx context.Context, sess *request.Session) error
	PingSession(ctx context.Context, ps *domain.PingSession) error
	GetOnlineDashboard(ctx context.Context) ([]domain.Session, error)
	IsOnlineSession(ctx context.Context, comp_name, login string) (bool, error)
}

func NewStorage(pool *pgxpool.Pool) Storage {
	return &storage{pool: pool}
}

type storage struct {
	pool *pgxpool.Pool
}

var (
	createSessionQuery = `INSERT INTO sessions (id, comp_name, ip_addr, login, next_ping_sec, date_time) 
	VALUES ($1, $2, $3, $4, $5, $6);`

	pingSessionQuery = `INSERT INTO sessions_ping (session_id, session_type, next_ping_date)
	VALUES ($1, $2, $3)
	ON CONFLICT (session_id, session_type)
	DO UPDATE SET
		next_ping_date = EXCLUDED.next_ping_date,
		updated = current_timestamp;`

	createOnlineSessionQuery = `INSERT INTO online_dashboard (session_id, comp_name, ip_addr, login, next_ping_date)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (comp_name)
	DO UPDATE SET
		session_id = EXCLUDED.session_id,
		ip_addr = EXCLUDED.ip_addr,
		login = EXCLUDED.login,
		next_ping_date = EXCLUDED.next_ping_date
		updated = current_timestamp;`

	onlineDashboardQuery = `SELECT session_id, comp_name, ip_addr, login, next_ping_date
	FROM online_dashboard;`

	isOnlineSessionQuery = `SELECT EXISTS (
		SELECT 1
		FROM online_dashboard
		WHERE comp_name = $1
		AND login = $2
		AND next_ping_date >= NOW()
	) AS is_online;`
)

func (s *storage) CreateSession(ctx context.Context, sess *request.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := s.pool.Exec(ctx, createSessionQuery,
		sess.ID,
		sess.ComputerName,
		sess.IPAddress,
		sess.Login,
		sess.NextPingSeconds,
		sess.DateTime,
	); err != nil {
		return err
	}

	// put in online dashboard
	if _, err := s.pool.Exec(ctx, createOnlineSessionQuery,
		sess.ID,
		sess.ComputerName,
		sess.IPAddress,
		sess.Login,
		sess.DateTime,
	); err != nil {
		return fmt.Errorf("online dashboard: %w", err)
	}

	return nil
}

func (s *storage) PingSession(ctx context.Context, ps *domain.PingSession) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := s.pool.Exec(ctx, pingSessionQuery,
		ps.SessionID,
		ps.SessionType,
		ps.NextPingDate,
	); err != nil {
		return err
	}

	// put in online dashboard
	if _, err := s.pool.Exec(ctx, createOnlineSessionQuery,
		ps.SessionID,
		ps.ComputerName,
		ps.IPAddress,
		ps.Login,
		ps.NextPingDate,
	); err != nil {
		return fmt.Errorf("online dashboard: %w", err)
	}

	return nil
}

func (s *storage) GetOnlineDashboard(ctx context.Context) ([]domain.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
			&s.NextPingDate,
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

func (s *storage) IsOnlineSession(ctx context.Context, comp_name, login string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var isOnline bool

	if err := s.pool.QueryRow(ctx, isOnlineSessionQuery,
		comp_name,
		login,
	).Scan(&isOnline); err != nil {
		return false, err
	}

	return isOnline, nil
}
