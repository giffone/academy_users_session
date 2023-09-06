package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"session_manager/internal/domain"
	"session_manager/internal/domain/request"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	CreateUsers(ctx context.Context, req []request.User) (err error)
	CreateComputers(ctx context.Context, req []request.Computer) (err error)
	CreateSession(ctx context.Context, sess *request.Session) error
	Activity(ctx context.Context, ps *request.Activity) error
	GetOnlineDashboard(ctx context.Context) ([]domain.Session, error)
	IsSessionExists(ctx context.Context, login string) (*domain.Session, error)
}

func NewStorage(pool *pgxpool.Pool) Storage {
	return &storage{pool: pool}
}

type storage struct {
	pool *pgxpool.Pool
}

const (
	// allow the "online" session notification to end (if late) or show that the user is no longer online
	// by subtracting n-seconds from the current time when checking the online session.
	//also affects how quickly a user can login again. so the interval should not be long (no more than 60 seconds).
	minusNSeconds = 10
)

func (s *storage) CreateUsers(ctx context.Context, req []request.User) (err error) {
	batch := &pgx.Batch{}

	for _, user := range req {
		if user.Name == "" {
			continue
		}
		batch.Queue(`INSERT INTO users (login) 
		VALUES ($1) 
		ON CONFLICT (login) DO NOTHING`,
			user.Name,
		)
	}

	ctx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	results := s.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < len(req); i++ {
		if _, err2 := results.Exec(); err2 != nil {
			err = errors.Join(err, err2)
		}
	}

	return err
}

func (s *storage) CreateComputers(ctx context.Context, req []request.Computer) (err error) {
	batch := &pgx.Batch{}

	for _, computer := range req {
		if computer.Name == "" {
			continue
		}
		batch.Queue(`INSERT INTO computers (comp_name) 
		VALUES ($1) 
		ON CONFLICT (comp_name) DO NOTHING`,
			computer.Name,
		)
	}

	ctx, cancel := context.WithTimeout(ctx, 180*time.Second)
	defer cancel()

	results := s.pool.SendBatch(ctx, batch)
	defer results.Close()

	for i := 0; i < len(req); i++ {
		if _, err2 := results.Exec(); err2 != nil {
			err = errors.Join(err, err2)
		}
	}

	return err
}

func (s *storage) CreateSession(ctx context.Context, req *request.Session) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// start session
	if _, err := s.pool.Exec(ctx, `INSERT INTO 
		sessions (id, comp_name, ip_addr, login, next_ping_sec, start_date_time, end_date_time) 
		VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		req.ID,
		req.ComputerName,
		req.IPAddress,
		req.Login,
		req.NextPingSeconds,
		req.DateTime,
		req.DateTime.Add(time.Duration(req.NextPingSeconds)*time.Second),
	); err != nil {
		return err
	}

	return nil
}

func (s *storage) Activity(ctx context.Context, req *request.Activity) error {
	ctx2, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	updateSessionEndQuery := `UPDATE sessions 
	SET end_date_time = $1
	WHERE id = $2;`

	end_date_time := req.DateTime.Add(time.Duration(req.NextPingSeconds) * time.Second)

	// -------------- if only session
	if req.SessionType == "" {
		if _, err := s.pool.Exec(ctx2, updateSessionEndQuery,
			end_date_time,
			req.SessionID,
		); err != nil {
			return fmt.Errorf("exec: activity: %w", err)
		}
		return nil
	}

	// -------------- if other activity [on zero platforn and etc...]
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err = tx.Rollback(ctx); err != nil {
			log.Printf("Activity: rollback: %s", err.Error())
		}
	}()

	// start activity
	if _, err := tx.Exec(ctx2,
		`INSERT INTO activity (session_id, session_type, login, start_date_time, end_date_time)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (session_id, session_type)
			DO UPDATE SET
			end_date_time = EXCLUDED.end_date_time;`,
		req.SessionID,
		req.SessionType,
		req.Login,
		req.DateTime,
		end_date_time,
	); err != nil {
		return fmt.Errorf("exec: activity: %w", err)
	}

	// also update session end_date_time
	if _, err := tx.Exec(ctx2, updateSessionEndQuery,
		end_date_time,
		req.SessionID,
	); err != nil {
		return fmt.Errorf("exec: activity: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (s *storage) GetOnlineDashboard(ctx context.Context) ([]domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx,
		`SELECT id, comp_name, ip_addr, login, start_date_time, end_date_time
		FROM sessions
		WHERE end_date_time >= NOW();`,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	sessions := make([]domain.Session, 0, 250)

	for rows.Next() {
		s := domain.Session{}
		if err := rows.Scan(
			&s.ID,
			&s.ComputerName,
			&s.IPAddress,
			&s.Login,
			&s.StartDateTime,
			&s.EndDateTime,
		); err != nil {
			return nil, fmt.Errorf("in iterate row: %w", err)
		}
		sessions = append(sessions, s)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows at all: %w", err)
	}

	return sessions, nil
}

func (s *storage) IsSessionExists(ctx context.Context, login string) (*domain.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var session domain.Session

	if err := s.pool.QueryRow(ctx,
		`SELECT id, comp_name, ip_addr, login, start_date_time, end_date_time
		FROM sessions
		WHERE login = $1
		AND end_date_time >= (NOW() - INTERVAL '$2 seconds');`,
		login,
		minusNSeconds,
	).Scan(&session); err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}
