package response

import (
	"time"
)

type Data struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type Session struct {
	ID            string    `db:"id" json:"id"`
	ComputerName  string    `db:"comp_name" json:"comp_name"`
	IPAddress     string    `db:"ip_addr" json:"ip_addr"`
	Login         string    `db:"login" json:"login"`
	StartDateTime time.Time `db:"start_date_time" json:"start_date_time"`
	EndDateTime   time.Time `db:"end_date_time" json:"end_date_time"`
}

type Activity struct {
	Login        string `db:"login"`
	TotalHours   int    `json:"total_hours"`
	UserActivity any    `json:"user_activity,omitempty"`
}

type UserActivityByMonth struct {
	Year  string `db:"year"`
	Month string `db:"month_name"`
	Hours int    `db:"hours"`
}

type UserActivityByDate struct {
	Date  time.Time `db:"date"`
	Hours int       `db:"hours"`
}
