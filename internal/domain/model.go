package domain

import "time"

type Session struct {
	ID            string    `db:"id"`
	ComputerName  string    `db:"comp_name"`
	IPAddress     string    `db:"ip_addr"`
	Login         string    `db:"login"`
	StartDateTime time.Time `db:"start_date_time"`
	EndDateTime   time.Time `db:"end_date_time"`
}