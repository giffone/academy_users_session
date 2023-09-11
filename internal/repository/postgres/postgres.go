package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"session_manager/internal/domain"
	"session_manager/internal/domain/request"
	"session_manager/internal/domain/response"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage interface {
	CreateUsers(ctx context.Context, req []request.User) (err error)
	CreateComputers(ctx context.Context, req []request.Computer) (err error)
	CreateSession(ctx context.Context, dto *domain.Session) error
	CreateActivity(ctx context.Context, dto *domain.Activity) error
	GetOnlineDashboard(ctx context.Context) ([]response.Session, error)
	IsSessionExists(ctx context.Context, login string) ([]response.Session, error)
	GetUserActivityByMonth(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error)
	GetUserActivityByDate(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error)
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
		batch.Queue(`INSERT INTO public.users (login) 
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
		batch.Queue(`INSERT INTO session.computers (comp_name) 
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

func (s *storage) CreateSession(ctx context.Context, dto *domain.Session) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// start session
	if _, err := s.pool.Exec(ctx, `INSERT INTO 
		session.in_campus (id, comp_name, ip_addr, login, next_ping_sec, start_date_time, end_date_time) 
		VALUES ($1, $2, $3, $4, $5, $6, $7);`,
		dto.ID,
		dto.ComputerName,
		dto.IPAddress,
		dto.Login,
		int(dto.NextPing.Seconds()),
		dto.StartDateTime,
		dto.EndDateTime,
	); err != nil {
		return customErr("exec", err)
	}

	return nil
}

func (s *storage) CreateActivity(ctx context.Context, dto *domain.Activity) error {
	ctx2, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	updateSessionEndQuery := `UPDATE session.in_campus
	SET end_date_time = $1
	WHERE id = $2;`

	// -------------- if only session
	if dto.SessionType == "" {
		if tag, err := s.pool.Exec(ctx2, updateSessionEndQuery,
			dto.EndDateTime,
			dto.SessionID,
		); err != nil {
			return customErr("session: exec: update", err)
		} else if tag.RowsAffected() == 0 {
			return &response.ErrNotFound
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
			log.Printf("CreateActivity: rollback: %s", err.Error())
		}
	}()

	// start activity
	if _, err := tx.Exec(ctx2,
		`INSERT INTO session.activity (session_id, session_type, login, start_date_time, end_date_time)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (session_id, session_type)
			DO UPDATE SET
			end_date_time = EXCLUDED.end_date_time;`,
		dto.SessionID,
		dto.SessionType,
		dto.Login,
		dto.StartDateTime,
		dto.EndDateTime,
	); err != nil {
		return customErr("activity: exec: insert", err)
	}

	// also update session end_date_time
	if tag, err := tx.Exec(ctx2, updateSessionEndQuery,
		dto.EndDateTime,
		dto.SessionID,
	); err != nil {
		return customErr("activity: exec: update", err)
	} else if tag.RowsAffected() == 0 {
		return &response.ErrNotFound
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func (s *storage) GetOnlineDashboard(ctx context.Context) ([]response.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	rows, err := s.pool.Query(ctx,
		`SELECT id, comp_name, ip_addr, login, start_date_time, end_date_time
		FROM session.in_campus
		WHERE end_date_time >= NOW();`,
	)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	sessions := make([]response.Session, 0, 250)

	for rows.Next() {
		session := response.Session{}
		if err := rows.Scan(
			&session.ID,
			&session.ComputerName,
			&session.IPAddress,
			&session.Login,
			&session.StartDateTime,
			&session.EndDateTime,
		); err != nil {
			return nil, fmt.Errorf("in iterate row: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows at all: %w", err)
	}

	return sessions, nil
}

func (s *storage) GetUserActivityByMonth(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var rows pgx.Rows
	var err error

	if dto.SessionType != "" {
		// from activity
		rows, err = s.pool.Query(ctx,
			`WITH monthly_hours AS (
				SELECT
					login,
					EXTRACT(YEAR FROM start_date_time) AS year,
					EXTRACT(MONTH FROM start_date_time) AS month_number,
					EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600 AS hours_calc
				FROM session.activity
				WHERE
					login = $1
					AND session_type = $2
					AND DATE_TRUNC('day', start_date_time) >= $3::date
					AND DATE_TRUNC('day', end_date_time) <= $4::date
			)
			SELECT 
				login,
				year,
				month_number,
				SUM(hours_calc) AS total_hours,
				SUM(SUM(hours_calc)) OVER (PARTITION BY login) AS total_hours
			FROM monthly_hours
			GROUP BY login, year, month_number
			ORDER BY year DESC, month_number DESC;`,
			dto.Login,
			dto.SessionType,
			dto.FromDate,
			dto.ToDate,
		)
	} else {
		// from in_campus
		rows, err = s.pool.Query(ctx,
			`WITH monthly_hours AS (
				SELECT
					login,
					EXTRACT(YEAR FROM start_date_time) AS year,
					EXTRACT(MONTH FROM start_date_time) AS month_number,
					EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600 AS hours_calc
				FROM session.in_campus
				WHERE
					login = $1
					AND DATE_TRUNC('day', start_date_time) >= $2::date
					AND DATE_TRUNC('day', end_date_time) <= $3::date
			)
			SELECT 
				login,
				year,
				month_number,
				SUM(hours_calc) AS total_hours,
				SUM(SUM(hours_calc)) OVER (PARTITION BY login) AS total_hours
			FROM monthly_hours
			GROUP BY login, year, month_number
			ORDER BY year DESC, month_number DESC;`,
			dto.Login,
			dto.FromDate,
			dto.ToDate,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	defer rows.Close()

	activities := make([]response.UserActivityByMonth, 0, 36)
	var totalHours float64
	var login string

	for rows.Next() {
		activity := response.UserActivityByMonth{}
		if err := rows.Scan(
			&login,
			&activity.Year,
			&activity.MonthNumber,
			&activity.Hours,
			&totalHours,
		); err != nil {
			return nil, fmt.Errorf("in iterate row: %w", err)
		}
		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows at all: %w", err)
	}

	return &response.UserActivity{
		Login:        login,
		TotalHours:   float32(math.Round(totalHours*100) / 100),
		UserActivity: activities,
	}, nil
}

func (s *storage) GetUserActivityByDate(ctx context.Context, dto *domain.UserActivity) (*response.UserActivity, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var rows pgx.Rows
	var err error
	if dto.SessionType != "" {
		// from activity
		rows, err = s.pool.Query(ctx,
			`SELECT
				login,
				DATE_TRUNC('day', start_date_time) AS date,
				SUM(EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600) AS hours,
				SUM(EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600) OVER (PARTITION BY login) AS total_hours
			FROM session.activity
			WHERE
				login = $1
				AND session_type = $2
				AND DATE_TRUNC('day', start_date_time) >= $3::date
				AND DATE_TRUNC('day', end_date_time) <= $4::date
			GROUP BY login, date, start_date_time, end_date_time
			ORDER BY date;`,
			dto.Login,
			dto.SessionType,
			dto.FromDate,
			dto.ToDate,
		)
	} else {
		// from in_campus
		rows, err = s.pool.Query(ctx,
			`SELECT
				login,
				DATE_TRUNC('day', start_date_time) AS date,
				SUM(EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600) AS hours,
				SUM(EXTRACT(EPOCH FROM (end_date_time - start_date_time)) / 3600) OVER (PARTITION BY login) AS total_hours
			FROM session.in_campus
			WHERE
				login = $1
				AND DATE_TRUNC('day', start_date_time) >= $2::date
				AND DATE_TRUNC('day', end_date_time) <= $3::date
			GROUP BY login, date, start_date_time, end_date_time
			ORDER BY date;`,
			dto.Login,
			dto.FromDate,
			dto.ToDate,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	activities := make([]response.UserActivityByDate, 0, 360)
	var totalHours float64
	var login string

	for rows.Next() {
		activity := response.UserActivityByDate{}
		if err := rows.Scan(
			&login,
			&activity.Date,
			&activity.Hours,
			&totalHours,
		); err != nil {
			return nil, fmt.Errorf("in iterate row: %w", err)
		}
		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows at all: %w", err)
	}

	return &response.UserActivity{
		Login:        login,
		TotalHours:   float32(math.Round(totalHours*100) / 100),
		UserActivity: activities,
	}, nil
}

func (s *storage) IsSessionExists(ctx context.Context, login string) ([]response.Session, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	query := fmt.Sprintf(`SELECT id, comp_name, ip_addr, login, start_date_time, end_date_time
	FROM session.in_campus
	WHERE login = $1
		AND end_date_time >= (NOW() - INTERVAL '%d seconds');`, minusNSeconds)

	rows, err := s.pool.Query(ctx, query, login)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	sessions := make([]response.Session, 0, 250)

	for rows.Next() {
		session := response.Session{}
		if err := rows.Scan(
			&session.ID,
			&session.ComputerName,
			&session.IPAddress,
			&session.Login,
			&session.StartDateTime,
			&session.EndDateTime,
		); err != nil {
			return nil, fmt.Errorf("in iterate row: %w", err)
		}
		sessions = append(sessions, session)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows at all: %w", err)
	}

	return sessions, nil
}

func customErr(message string, err error) error {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == pgerrcode.UniqueViolation {
			return &response.ErrDuplicateKey
		}
		if pgErr.Code == pgerrcode.RaiseException {
			if pgErr.Message == response.ErrEndStartDate.Error() {
				return &response.ErrEndStartDate
			}
			if pgErr.Message == response.ErrEndEndDate.Error() {
				return &response.ErrEndEndDate
			}
		}
	}
	if message != "" {
		return fmt.Errorf("%s: %s", message, err)
	}
	return err
}
