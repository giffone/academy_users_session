package request

import (
	"errors"
	"time"
)

type User struct {
	Name string `json:"name"`
}

type Computer struct {
	Name string `json:"name"`
}

type Session struct {
	ID              string    `json:"id"`
	ComputerName    string    `json:"comp_name"`
	IPAddress       string    `json:"ip_addr"`
	Login           string    `json:"login"`
	NextPingSeconds int       `json:"next_ping_sec"`
	DateTime        time.Time `json:"date_time"`
}

func (s *Session) Validate() error {
	if s.ID == "" {
		return errors.New("id is empty")
	}
	if s.ComputerName == "" {
		return errors.New("comp_name is empty")
	}
	// if s.IPAddress == "" {
	// 	return response.ErrEmpty("ip")
	// }
	if s.Login == "" {
		return errors.New("login is empty")
	}
	if s.NextPingSeconds <= 0 {
		return errors.New("next ping duration less or eq 0")
	}
	if s.DateTime.IsZero() {
		s.DateTime = time.Now()
	}
	return nil
}

type Activity struct {
	SessionID       string    `json:"session_id"`
	SessionType     string    `json:"session_type"`
	Login           string    `json:"login"`
	NextPingSeconds int       `json:"next_ping_sec"`
	DateTime        time.Time `json:"date_time"`
}

func (a *Activity) Validate() error {
	if a.SessionID == "" {
		return errors.New("session_id is empty")
	}
	if a.Login == "" {
		return errors.New("login is empty")
	}
	if a.NextPingSeconds <= 0 {
		return errors.New("next ping duration less or eq 0")
	}
	if a.DateTime.IsZero() {
		a.DateTime = time.Now()
	}
	return nil
}

type UserActivity struct {
	SessionType string    `query:"session_type"` // parsing by link's queries
	Login       string    `query:"login"`
	FromDate    time.Time `query:"from_date"`
	ToDate      time.Time `query:"to_date"`
	GroupBy     string    `query:"group_by"`
}

const (
	GroupByMonth = "month"
	GroupByAate  = "date"
)

func (ua *UserActivity) Validate() error {
	if ua.SessionType == "" {
		return errors.New("session_type is empty")
	}
	if ua.Login == "" {
		return errors.New("login is empty")
	}
	if ua.FromDate.IsZero() && ua.ToDate.IsZero() {
		ua.FromDate = time.Now().Truncate(24 * time.Hour)
		ua.ToDate = ua.FromDate.Add(24 * time.Hour)
	}
	if ua.ToDate.IsZero() {
		ua.FromDate = time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)
	}
	if ua.GroupBy == "" {
		ua.GroupBy = GroupByAate
	}
	if ua.GroupBy != GroupByMonth && ua.GroupBy != GroupByAate {
		return errors.New("group by must be 'month' or 'date'")
	}
	return nil
}
