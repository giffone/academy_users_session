package postgres

import "session_manager/internal/domain"

type Storage interface {
	CreateSession(sess *domain.Session) error
	UpdateSession(sess *domain.Session) error
	GetOnlineSessions() (domain.Sessions, error)
}

func New() Storage {
	return &storage{}
}

type storage struct {
}


func (s *storage)  CreateSession(sess *domain.Session) error {
	return nil
}

func (s *storage)  UpdateSession(sess *domain.Session) error {
	return nil
}

func (s *storage)  GetOnlineSessions() (domain.Sessions, error) {
	return nil, nil
}