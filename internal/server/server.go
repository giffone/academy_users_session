package server

import (
	"session_manager/internal/api"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
}

func New() {

}

func newRouter(h api.Handlers) *echo.Echo {
	e := echo.New()

	// set middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// register handlers
	g := e.Group("/api/session-manager")
	g.POST("/session", h.CreateSession)
	g.PUT("/session/:id", h.UpdateSession)
	g.GET("/online-sessions", h.GetOnlineSessions)
	return e
}
