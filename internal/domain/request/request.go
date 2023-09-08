package request

import (
	"errors"
	"session_manager/internal/domain"
	"time"
)

type User struct {
	Name string `json:"name"`
}

type Computer struct {
	Name string `json:"name"`
}

type Session struct {
	ID              string `json:"id"`
	ComputerName    string `json:"comp_name"`
	IPAddress       string `json:"ip_addr"`
	Login           string `json:"login"`
	NextPingSeconds int    `json:"next_ping_sec"`
	DateTime        string `json:"date_time"`
}

func (s *Session) Validate() (*domain.Session, error) {
	if s.ID == "" {
		return nil, errors.New("id is empty")
	}
	if s.ComputerName == "" {
		return nil, errors.New("comp_name is empty")
	}
	if s.Login == "" {
		return nil, errors.New("login is empty")
	}
	if s.NextPingSeconds <= 0 {
		return nil, errors.New("next ping duration less or eq 0")
	}
	dto := domain.Session{
		ID:           s.ID,
		ComputerName: s.ComputerName,
		IPAddress:    s.IPAddress,
		Login:        s.Login,
		NextPing:     time.Duration(s.NextPingSeconds) * time.Second,
	}
	if s.DateTime == "" {
		dto.StartDateTime = time.Now()
	} else {
		t, err := parseDate(s.DateTime)
		if err != nil {
			return nil, err
		}
		dto.StartDateTime = t
	}
	dto.EndDateTime = dto.StartDateTime.Add(dto.NextPing)
	return &dto, nil
}

type Activity struct {
	SessionID       string `json:"session_id"`
	SessionType     string `json:"session_type,omitempty"`
	Login           string `json:"login"`
	NextPingSeconds int    `json:"next_ping_sec"`
	DateTime        string `json:"date_time"`
}

func (a *Activity) Validate() (*domain.Activity, error) {
	if a.SessionID == "" {
		return nil, errors.New("session_id is empty")
	}
	if a.Login == "" {
		return nil, errors.New("login is empty")
	}
	if a.NextPingSeconds <= 0 {
		return nil, errors.New("next ping duration less or eq 0")
	}
	dto := domain.Activity{
		SessionID:   a.SessionID,
		SessionType: a.SessionType,
		Login:       a.Login,
	}
	if a.DateTime == "" {
		dto.StartDateTime = time.Now()
	} else {
		t, err := parseDate(a.DateTime)
		if err != nil {
			return nil, err
		}
		dto.StartDateTime = t
	}
	dto.EndDateTime = dto.StartDateTime.Add(time.Duration(a.NextPingSeconds) * time.Second)
	return &dto, nil
}

type UserActivity struct {
	SessionType string `query:"session_type"` // parsing by link's queries ('omitempty' not working, do not add)
	Login       string `query:"login"`
	FromDate    string `query:"from_date"`
	ToDate      string `query:"to_date"`
	GroupBy     string `query:"group_by"`
}

const (
	GroupByMonth = "month"
	GroupByDate  = "date"
)

func (ua *UserActivity) Validate() (*domain.UserActivity, error) {
	if ua.Login == "" {
		return nil, errors.New("login is empty")
	}
	if ua.GroupBy == "" {
		ua.GroupBy = GroupByDate
	}
	if ua.GroupBy != GroupByMonth && ua.GroupBy != GroupByDate {
		return nil, errors.New("group by must be 'month' or 'date'")
	}

	dto := domain.UserActivity{
		SessionType: ua.SessionType,
		Login:       ua.Login,
		GroupBy:     ua.GroupBy,
	}

	if ua.FromDate == "" && ua.ToDate == "" {
		dto.FromDate = time.Now().Truncate(24 * time.Hour)
		dto.ToDate = dto.FromDate.Add(24 * time.Hour)
		return &dto, nil
	}

	t, err := parseDate(ua.FromDate)
	if err != nil {
		return nil, err
	}
	dto.FromDate = t.Truncate(24 * time.Hour)

	if ua.ToDate == "" {
		dto.ToDate = time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)
	} else {
		t, err := parseDate(ua.ToDate)
		if err != nil {
			return nil, err
		}
		dto.ToDate = t.Truncate(24 * time.Hour)
	}

	return &dto, nil
}

func parseDate(s string) (t time.Time, err error) {
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse(time.DateTime, s)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse(time.DateOnly, s)
	if err == nil {
		return t, nil
	}
	return t, errors.New(`date format must be 
	'2006-01-02T15:04:05Z07:00' or 
	'2006-01-02 15:04:05' or 
	'2006-01-02'`,
	)
}
