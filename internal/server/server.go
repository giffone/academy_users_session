package server

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"session_manager/internal/api"
	"session_manager/internal/repository/postgres"
	"session_manager/internal/service"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server interface {
	Run(ctx context.Context)
	Stop(ctx context.Context)
}

type server struct {
	router   *echo.Echo
}

func NewServer(env *Env) Server {
	s := server{
		router: echo.New(),
	}

	// storage
	storage := postgres.NewStorage(env.pool)

	// service
	svc := service.New(storage)

	// handlers
	hndl := api.NewHandlers(s.router.Logger, svc)

	// set middlewares
	s.router.Use(middleware.Logger(), middleware.Recover())

	// register handlers
	g := s.router.Group("/api/session-manager")
	g.POST("/session", hndl.CreateSession)
	g.POST("/activity", hndl.Activity)
	g.GET("/dashboard", hndl.GetOnlineSessions)

	return &s
}

func (s *server) Run(ctx context.Context) {
	ctxSignal, cancelSignal := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// start server
	go func() {
		if err := s.router.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Printf("server start error: %s\n", err.Error())
			cancelSignal()
		}
	}()

	// wait system notifiers or cancel func
	<-ctxSignal.Done()
}

func (s *server) Stop(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := s.router.Shutdown(ctx); err != nil {
		log.Printf("server stop error: %s\n", err.Error())
		return
	}
	log.Println("server stopped successfully")
}