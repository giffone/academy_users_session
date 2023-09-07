package service

import (
	"context"
	"fmt"
	"session_manager/internal/domain/request"
	"session_manager/internal/domain/response"
	"session_manager/internal/repository/postgres"
)

type Service interface {
	CreateUsers(ctx context.Context, req []request.User) error
	CreateComputers(ctx context.Context, req []request.Computer) error
	CreateSession(ctx context.Context, req *request.Session) ([]response.Session, error)
	CreateActivity(ctx context.Context, req *request.Activity) error
	GetOnlineDashboard(ctx context.Context) ([]response.Session, error)
	GetUserActivity(ctx context.Context, req *request.UserActivity) (activity *response.Activity, err error)
}

func New(storage postgres.Storage) Service {
	return &service{storage: storage}
}

type service struct {
	storage postgres.Storage
}

func (s *service) CreateUsers(ctx context.Context, req []request.User) error {
	return s.storage.CreateUsers(ctx, req)
}

func (s *service) CreateComputers(ctx context.Context, req []request.Computer) error {
	return s.storage.CreateComputers(ctx, req)
}

func (s *service) CreateSession(ctx context.Context, req *request.Session) ([]response.Session, error) {
	// first check if session already exists
	if sessions, err := s.storage.IsSessionExists(ctx, req.Login); err != nil {
		return nil, fmt.Errorf("IsSessionExists: %w", err)
	} else if len(sessions) != 0 {
		return sessions, response.ErrAccessDenied
	}

	// create session
	return nil, s.storage.CreateSession(ctx, req)
}

func (s *service) CreateActivity(ctx context.Context, req *request.Activity) error {
	return s.storage.CreateActivity(ctx, req)
}

func (s *service) GetOnlineDashboard(ctx context.Context) ([]response.Session, error) {
	return s.storage.GetOnlineDashboard(ctx)
}

func (s *service) GetUserActivity(ctx context.Context, req *request.UserActivity) (activity *response.Activity, err error) {
	if req.GroupBy == request.GroupByMonth {
		return s.storage.GetUserActivityByMonth(ctx, req)
	}

	return s.storage.GetUserActivityByDate(ctx, req)
}
