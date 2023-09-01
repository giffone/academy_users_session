package domain

type Session struct {
	ID           string `json:"id" db:"id"`
	ComputerName string `json:"comp_name" db:"comp_name"`
	IPAddress    string `json:"ip_addr" db:"ip_addr"`
	Login        string `json:"login" db:"login"`
	Status       string `json:"status" db:"status"`
	DateTime     int    `json:"date_time" db:"date_time"`
}

func (s *Session) validate() *Response {
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
	if s.Status == "" {
		return errEmpty("status")
	}
	if s.DateTime == 0 {
		return errEmpty("date_time")
	}
	return nil
}
