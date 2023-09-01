package domain

import "fmt"

type Response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
}

func errEmpty(field string) *Response {
	return &Response{
		Code:   400,
		Status: fmt.Sprintf("%s is empty", field),
	}
}
