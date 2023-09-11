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

type UserActivity struct {
	Login        string  `db:"login" json:"id"`
	TotalHours   float32 `db:"total_hours" json:"total_hours"`
	UserActivity any     `json:"user_activity,omitempty"`
}

type UserActivityByMonth struct {
	Year        string  `db:"year" json:"year"`
	MonthNumber string  `db:"month_number" json:"month_num"`
	Hours       float32 `db:"hours" json:"hours"`
}

type UserActivityByDate struct {
	Date  time.Time `db:"date" json:"date"`
	Hours float32   `db:"hours" json:"hours"`
}
