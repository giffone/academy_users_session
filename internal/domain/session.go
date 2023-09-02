package domain

import "time"

type Session struct {
	ID           string    `json:"id" db:"id"`
	ComputerName string    `json:"comp_name" db:"comp_name"`
	IPAddress    string    `json:"ip_addr" db:"ip_addr"`
	Login        string    `json:"login" db:"login"`
	DateTime     time.Time `json:"date_time" db:"date_time"`
}

func (s *Session) Validate() *Response {
	if s.ID == "" {
		return errEmpty("id")
	}
	if s.ComputerName == "" {
		return errEmpty("comp_name")
	}
	if s.IPAddress == "" {
		return errEmpty("ip")
	}
	if s.Login == "" {
		return errEmpty("login")
	}
	if s.DateTime.IsZero() {
		return errEmpty("date_time")
	}
	return nil
}

type PingSession struct {
	SessionID    string    `json:"session_id" db:"session_id"`
	SessionType  string    `json:"session_type" db:"session_type"`
	ComputerName string    `json:"comp_name" db:"comp_name"`
	IPAddress    string    `json:"ip_addr" db:"ip_addr"`
	Login        string    `json:"login" db:"login"`
	DateTime     time.Time `json:"date_time" db:"date_time"`
}

func (ps *PingSession) Validate() *Response {
	if ps.SessionID == "" {
		return errEmpty("session_id")
	}
	if ps.SessionType == "" {
		return errEmpty("session_type")
	}
	if ps.ComputerName == "" {
		return errEmpty("comp_name")
	}
	if ps.IPAddress == "" {
		return errEmpty("ip")
	}
	if ps.Login == "" {
		return errEmpty("login")
	}
	if ps.DateTime.IsZero() {
		return errEmpty("date_time")
	}
	return nil
}

type Sessions []Session
