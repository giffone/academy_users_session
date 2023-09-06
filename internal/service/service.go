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
	CreateSession(ctx context.Context, req *request.Session) (*domain.Session, error)
	Activity(ctx context.Context, req *request.Activity) error
	GetOnlineSessions(ctx context.Context) ([]domain.Session, error)
}

func New(storage postgres.Storage) Service {
	return &service{storage: storage}
}

type service struct {
	storage postgres.Storage
}

func (s *service) CreateUsers(ctx context.Context, req []request.User) error {
	if err := s.storage.CreateUsers(ctx, req); err != nil {
		return fmt.Errorf("CreateUsers: %w", err)
	}

	return nil
}

func (s *service) CreateComputers(ctx context.Context, req []request.Computer) error {
	if err := s.storage.CreateComputers(ctx, req); err != nil {
		return fmt.Errorf("CreateComputers: %w", err)
	}

	return nil
}

func (s *service) CreateSession(ctx context.Context, req *request.Session) (*domain.Session, error) {
	// first check if session already exists
	if session, err := s.storage.IsSessionExists(ctx, req.Login); err != nil {
		return nil, fmt.Errorf("IsSessionExists: %w", err)
	} else if session != nil {
		return session, response.ErrAccessDenied
	}

	// create session
	if err := s.storage.CreateSession(ctx, req); err != nil {
		return nil, fmt.Errorf("CreateSession: %w", err)
	}

	return nil, nil
}

func (s *service) Activity(ctx context.Context, req *request.Activity) error {
	if err := s.storage.Activity(ctx, req); err != nil {
		return fmt.Errorf("Activity: %w", err)
	}

	return nil
}

func (s *service) GetOnlineSessions(ctx context.Context) ([]domain.Session, error) {
	sessions, err := s.storage.GetOnlineDashboard(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetOnlineDashboard: %w", err)
	}
	return sessions, nil
}
