package api

import "github.com/labstack/echo/v4"

type Handlers interface {
	CreateSession(c echo.Context) error 
	UpdateSession(c echo.Context) error
	GetOnlineSessions(c echo.Context) error 
}

type handlers struct {
}

func NewHandlers() Handlers{
	return &handlers{}
}

func(h *handlers) CreateSession(c echo.Context) error {
	return nil
}

func(h *handlers)  UpdateSession(c echo.Context) error {
	return nil
}

func(h *handlers)  GetOnlineSessions(c echo.Context) error {
	return nil
}

