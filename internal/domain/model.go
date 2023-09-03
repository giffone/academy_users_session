package domain

import "time"

type Session struct {
	SessionID    string    `db:"session_id"`
	ComputerName string    `db:"comp_name"`
	IPAddress    string    `db:"ip_addr"`
	Login        string    `db:"login"`
	DateTime     time.Time `db:"date_time"`
}

// type PingSession struct {
// 	SessionID    string    `db:"session_id"`
// 	SessionType  string    `db:"session_type"`
// 	ComputerName string    `db:"comp_name"`
// 	IPAddress    string    `db:"ip_addr"`
// 	Login        string    `db:"login"`
// 	NextPingDate time.Time `db:"next_ping_date"`
// }
