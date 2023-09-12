package service

import (
	"context"
	"fmt"
	"session_manager/internal/domain"
	"session_manager/internal/domain/request"
	"session_manager/internal/domain/response"
	"session_manager/internal/repository/postgres"
)

type Service interface {
	CreateUsers(ctx context.Context, req []request.User) error
	CreateComputers(ctx context.Context, req []request.Computer) error
	CreateSession(ctx context.Context, dto *domain.Session) ([]response.Session, error)
	CreateActivity(ctx context.Context, dto *domain.Activity) error
	GetOnlineDashboard(ctx context.Context) ([]response.Session, error)
	GetUserActivity(ctx context.Context, dto *domain.UserActivity) (activity *response.UserActivity, err error)
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

func (s *service) CreateSession(ctx context.Context, dto *domain.Session) ([]response.Session, error) {
	// first check if session already exists
	if sessions, err := s.storage.IsSessionExists(ctx, dto.Login); err != nil {
		return nil, fmt.Errorf("IsSessionExists: %w", err)
	} else if len(sessions) != 0 {
		return sessions, response.ErrAccessDenied
	}

	// create session
	return nil, s.storage.CreateSession(ctx, dto)
}

func (s *service) CreateActivity(ctx context.Context, dto *domain.Activity) error {
	return s.storage.CreateActivity(ctx, dto)
}

func (s *service) GetOnlineDashboard(ctx context.Context) ([]response.Session, error) {
	return s.storage.GetOnlineDashboard(ctx)
}

func (s *service) GetUserActivity(ctx context.Context, dto *domain.UserActivity) (activity *response.UserActivity, err error) {
	if dto.GroupBy == request.GroupByMonth {
		if dto.SessionType == "" {
			// no need sort - get from main table in_campus
			return s.storage.GetUserActivityByMonthInCampus(ctx, dto)
		}
		// need sort
		return s.storage.GetUserActivityByMonth(ctx, dto)
	}

	if dto.SessionType == "" {
		// no need sort - get from main table in_campus
		return s.storage.GetUserActivityByDateInCampus(ctx, dto)
	}
	// need sort
	return s.storage.GetUserActivityByDate(ctx, dto)
}