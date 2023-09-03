package response

import (
	"fmt"
	"session_manager/internal/domain"
)

type Data struct {
	Code     int              `json:"code"`
	Status   string           `json:"status"`
	Sessions []domain.Session `json:"sessions,omitempty"`
}

func ErrEmpty(field string) *Data {
	return &Data{
		Code:   400,
		Status: fmt.Sprintf("%s is empty", field),
	}
}
