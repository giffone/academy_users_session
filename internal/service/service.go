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
	CreateSession(ctx context.Context, req *request.Session) (*domain.Session, error)
	Activity(ctx context.Context, req *request.Activity) error
	GetOnlineSessions() ([]domain.Session, error)
}

func New() Service {
	return &service{}
}

type service struct {
	storage postgres.Storage
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

func (s *service) GetOnlineSessions() ([]domain.Session, error) {
	// s.storage.GetOnlineDashboard()
	return nil, nil
}
