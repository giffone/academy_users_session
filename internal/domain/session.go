package domain

type Session struct {
	ID           string `json:"id" db:"id"`
	ComputerName string `json:"comp_name" db:"comp_name"`
	IPAddress    string `json:"ip_addr" db:"ip_addr"`
	Login        string `json:"login" db:"login"`
	Status       string `json:"status" db:"status"`
	DateTime     int    `json:"date_time" db:"date_time"`
}

func (s *Session) validate() error {
	return nil
}