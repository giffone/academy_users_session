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
	PingSession(ctx context.Context, req *request.PingSession) error
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
	if session, err := s.storage.IsSessionExists(ctx, req.ComputerName, req.Login); err != nil {
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

func (s *service) PingSession(ctx context.Context, req *request.PingSession) error {
	if err := s.storage.PingSession(ctx, req); err != nil {
		return fmt.Errorf("PingSession: %w", err)
	}

	return nil
}

func (s *service) GetOnlineSessions() ([]domain.Session, error) {
	// s.storage.GetOnlineDashboard()
	return nil, nil
}
