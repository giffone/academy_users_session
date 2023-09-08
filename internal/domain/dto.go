package domain

import (
	"time"
)

type Session struct {
	ID            string
	ComputerName  string
	IPAddress     string
	Login         string
	NextPing      time.Duration
	StartDateTime time.Time
	EndDateTime   time.Time
}

type Activity struct {
	SessionID     string
	SessionType   string
	Login         string
	StartDateTime time.Time
	EndDateTime   time.Time
}

type UserActivity struct {
	SessionType string
	Login       string
	FromDate    time.Time
	ToDate      time.Time
	GroupBy     string
}
