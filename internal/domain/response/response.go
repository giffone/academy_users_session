package response

import (
	"errors"
	"fmt"
	"session_manager/internal/domain"
)

type Data struct {
	Message  string           `json:"status"`
	Sessions []domain.Session `json:"sessions,omitempty"`
}

func ErrEmpty(field string) *Data {
	return &Data{
		Message: fmt.Sprintf("%s is empty", field),
	}
}

var ErrAccessDenied = errors.New("access denied")
