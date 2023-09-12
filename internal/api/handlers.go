package api

import (
	"errors"
	"fmt"
	"log"

	"net/http"
	"session_manager/internal/domain/request"
	"session_manager/internal/domain/response"
	"session_manager/internal/service"

	"github.com/labstack/echo/v4"
)

var logErr = "logErr"

type Handlers interface {
	CreateUsers(c echo.Context) error
	CreateComputers(c echo.Context) error
	CreateSession(c echo.Context) error
	CreateActivity(c echo.Context) error
	GetOnlineSessions(c echo.Context) error
	GetUserActivity(c echo.Context) error
}

type handlers struct {
	svc service.Service
}

func NewHandlers(logg echo.Logger, svc service.Service) Handlers {
	return &handlers{
		svc: svc,
	}
}

func (h *handlers) CreateUsers(c echo.Context) error {
	var req []request.User

	// parse data
	if err := c.Bind(&req); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateUsers: bind req body: %s", err))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// create users
	if err := h.svc.CreateUsers(c.Request().Context(), req); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateUsers: %s", err))
		return customErrResponse(c, err, nil)
	}

	return created(c)
}

func (h *handlers) CreateComputers(c echo.Context) error {
	var req []request.Computer

	// parse data
	if err := c.Bind(&req); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateComputers: bind req body: %s", err))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// create users
	if err := h.svc.CreateComputers(c.Request().Context(), req); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateComputers: %s", err))
		return customErrResponse(c, err, nil)
	}

	return created(c)
}

func (h *handlers) CreateSession(c echo.Context) error {
	var req request.Session

	// parse data
	if err := c.Bind(&req); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateSession: bind req body: %s", err))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// validate data
	dto, err := req.Validate()
	if err != nil {
		c.Set(logErr, fmt.Sprintf("CreateSession: validate: %s", err))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// create session
	if sess, err := h.svc.CreateSession(c.Request().Context(), dto); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateSession: %s", err))
		return customErrResponse(c, err, sess)
	}

	return created(c)
}

func (h *handlers) CreateActivity(c echo.Context) error {
	var req request.Activity

	// parse data
	if err := c.Bind(&req); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateActivity: bind req body: %s", err))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// validate data
	dto, err := req.Validate()
	if err != nil {
		c.Set(logErr, fmt.Sprintf("CreateActivity: validate: %s", err))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// create activity
	if err := h.svc.CreateActivity(c.Request().Context(), dto); err != nil {
		c.Set(logErr, fmt.Sprintf("CreateActivity: %s", err))
		return customErrResponse(c, err, nil)
	}

	return created(c)
}

func (h *handlers) GetOnlineSessions(c echo.Context) error {
	sessions, err := h.svc.GetOnlineDashboard(c.Request().Context())
	if err != nil {
		c.Set(logErr, fmt.Sprintf("GetOnlineSessions: %s", err))
		return customErrResponse(c, err, nil)
	}

	return ok(c, sessions)
}

func (h *handlers) GetUserActivity(c echo.Context) (err error) {
	var req request.UserActivity

	// parse data
	if err := c.Bind(&req); err != nil {
		c.Set(logErr, fmt.Sprintf("GetUserActivity: bind req body: %s", err))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	// validate data
	dto, err := req.Validate()
	if err != nil {
		c.Set(logErr, fmt.Sprintf("GetUserActivity: validate: %s", err))
		return customErrResponse(c, &response.ErrBadReq{Message: err.Error()}, nil)
	}

	activity, err := h.svc.GetUserActivity(c.Request().Context(), dto)
	if err != nil {
		c.Set(logErr, fmt.Sprintf("GetUserActivity: %s", err))
		return customErrResponse(c, err, nil)
	}

	return ok(c, activity)
}

func customErrResponse(c echo.Context, err error, data any) error {
	defer printLogErr(c)

	if data == nil {
		data = []string{} // to show empty array
	}
	if errors.Is(err, response.ErrAccessDenied) {
		return c.JSON(http.StatusUnauthorized, response.Data{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
			Data:    data,
		})
	}
	var errBadReq *response.ErrBadReq
	if errors.As(err, &errBadReq) {
		return c.JSON(http.StatusBadRequest, response.Data{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
			Data:    data,
		})
	}

	return c.JSON(http.StatusInternalServerError, response.Data{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
		Data:    data,
	})
}

func printLogErr(c echo.Context) {
	if eLog := c.Get(logErr); eLog != nil {
		log.Printf("\n//----\n[error]: %v\n----\\\\\n", eLog)
	}
}

func created(c echo.Context) error {
	return c.JSON(http.StatusCreated, response.Data{
		Code:    http.StatusCreated,
		Message: http.StatusText(http.StatusCreated),
		Data:    []string{},
	})
}

func ok(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, response.Data{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	})
}
