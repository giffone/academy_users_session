package api

import (
	"errors"
	"fmt"
	"net/http"
	"session_manager/internal/domain"
	"session_manager/internal/domain/request"
	"session_manager/internal/domain/response"
	"session_manager/internal/service"

	"github.com/labstack/echo/v4"
)

type Handlers interface {
	CreateSession(c echo.Context) error
	UpdateSession(c echo.Context) error
	GetOnlineSessions(c echo.Context) error
}

type handlers struct {
	logg echo.Logger
	svc  service.Service
}

func NewHandlers(logg echo.Logger, svc service.Service) Handlers {
	return &handlers{
		logg: logg,
		svc:  svc,
	}
}

func (h *handlers) CreateSession(c echo.Context) error {
	var req request.Session

	// parse data
	if err := c.Bind(&req); err != nil {
		e := fmt.Errorf("bind req body: %s", err)
		h.logg.Error(e)
		return c.JSON(http.StatusBadRequest, response.Data{Message: e.Error()})
	}

	// validate data
	if res := req.Validate(); res != nil {
		h.logg.Errorf("validate: %s", res.Message)
		return c.JSON(http.StatusBadRequest, res)
	}

	// create session
	if sess, err := h.svc.CreateSession(c.Request().Context(), &req); err != nil {
		if errors.Is(err, response.ErrAccessDenied) && sess != nil {
			return c.JSON(http.StatusUnauthorized, response.Data{
				Message:  err.Error(),
				Sessions: []domain.Session{*sess},
			})
		}
		h.logg.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, response.Data{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, nil)
}

func (h *handlers) UpdateSession(c echo.Context) error {
	return nil
}

func (h *handlers) GetOnlineSessions(c echo.Context) error {
	return nil
}
