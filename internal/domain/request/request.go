package request

import (
	"session_manager/internal/domain/response"
	"time"
)

type Session struct {
	ID              string    `json:"id"`
	ComputerName    string    `json:"comp_name"`
	IPAddress       string    `json:"ip_addr"`
	Login           string    `json:"login"`
	NextPingSeconds int       `json:"next_ping_sec"`
	DateTime        time.Time `json:"date_time"`
}

func (s *Session) Validate() *response.Data {
	if s.ID == "" {
		return response.ErrEmpty("id")
	}
	if s.ComputerName == "" {
		return response.ErrEmpty("comp_name")
	}
	// if s.IPAddress == "" {
	// 	return response.ErrEmpty("ip")
	// }
	if s.Login == "" {
		return response.ErrEmpty("login")
	}
	if s.NextPingSeconds <= 0 {
		return response.ErrEmpty("next ping duration less or eq 0")
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

func (a *Activity) Validate() *response.Data {
	if a.SessionID == "" {
		return response.ErrEmpty("session_id")
	}
	if a.SessionType == "" {
		return response.ErrEmpty("session_type")
	}
	if a.Login == "" {
		return response.ErrEmpty("login")
	}
	if a.NextPingSeconds <= 0 {
		return response.ErrEmpty("next ping duration less or eq 0")
	}
	if a.DateTime.IsZero() {
		a.DateTime = time.Now()
	}
	return nil
}
