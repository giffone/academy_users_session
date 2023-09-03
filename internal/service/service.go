package service

import (
	"session_manager/internal/domain"
	"session_manager/internal/domain/request"
)

type Service interface {
	CreateSession(sess *request.Session) error
	UpdateSession(sess *request.Session) error
	GetOnlineSessions() ([]domain.Session, error)
}

func New() Service {
	return &service{}
}

type service struct {
}

func (s *service) CreateSession(sess *request.Session) error {
	return nil
}

func (s *service) UpdateSession(sess *request.Session) error {
	return nil
}

func (s *service) GetOnlineSessions() ([]domain.Session, error) {
	return nil, nil
}