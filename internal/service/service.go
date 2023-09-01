package service

import "session_manager/internal/domain"

type Service interface {
	CreateSession(sess *domain.Session) error
	UpdateSession(sess *domain.Session) error
	GetOnlineSessions() (domain.Sessions, error)
}

func New() Service {
	return &service{}
}

type service struct {
}

func (s *service) CreateSession(sess *domain.Session) error {
	return nil
}

func (s *service) UpdateSession(sess *domain.Session) error {
	return nil
}

func (s *service) GetOnlineSessions() (domain.Sessions, error) {
	return nil, nil
}