package response

import "errors"

type ErrBadReq struct{ Message string }

func (e *ErrBadReq) Error() string { return e.Message }

var (
	ErrAccessDenied = errors.New("access denied")
	
	ErrDuplicateKey = &ErrBadReq{"duplicate key error"}
	ErrForeignKey   = &ErrBadReq{"unknown foreign key error"}
	ErrNotFound     = &ErrBadReq{"not found"}
	ErrEndStartDate = &ErrBadReq{"end_date_time must be greater than start_date_time"}
	ErrEndEndDate   = &ErrBadReq{"end_date_time must be greater than previous value"}
)
