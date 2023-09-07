package api

import (
	"errors"
	"fmt"
	"net/http"
	"session_manager/internal/domain/request"
	"session_manager/internal/domain/response"
	"session_manager/internal/service"

	"github.com/labstack/echo/v4"
)

type Handlers interface {
	CreateUsers(c echo.Context) error
	CreateComputers(c echo.Context) error
	CreateSession(c echo.Context) error
	CreateActivity(c echo.Context) error
	GetOnlineSessions(c echo.Context) error
	GetUserActivity(c echo.Context) error
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

func (h *handlers) CreateUsers(c echo.Context) error {
	var req []request.User

	// parse data
	if err := c.Bind(&req); err != nil {
		e := fmt.Errorf("bind req body: %s", err)
		h.logg.Error(e)
		return c.JSON(http.StatusBadRequest, response.Data{Message: e.Error()})
	}

	// create users
	if err := h.svc.CreateUsers(c.Request().Context(), req); err != nil {
		h.logg.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, response.Data{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, response.Data{
		Message: http.StatusText(http.StatusCreated)},
	)
}

func (h *handlers) CreateComputers(c echo.Context) error {
	var req []request.Computer

	// parse data
	if err := c.Bind(&req); err != nil {
		e := fmt.Errorf("bind req body: %s", err)
		h.logg.Error(e)
		return c.JSON(http.StatusBadRequest, response.Data{Message: e.Error()})
	}

	// create users
	if err := h.svc.CreateComputers(c.Request().Context(), req); err != nil {
		h.logg.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, response.Data{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, response.Data{
		Message: http.StatusText(http.StatusCreated)},
	)
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
				Message: err.Error(),
				Data:    []response.Session{*sess},
			})
		}
		h.logg.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, response.Data{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, response.Data{
		Message: http.StatusText(http.StatusCreated)},
	)
}

func (h *handlers) CreateActivity(c echo.Context) error {
	var req request.Activity

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

	// create activity
	if err := h.svc.CreateActivity(c.Request().Context(), &req); err != nil {
		h.logg.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, response.Data{Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, response.Data{
		Message: http.StatusText(http.StatusCreated)},
	)
}

func (h *handlers) GetOnlineSessions(c echo.Context) error {
	sessions, err := h.svc.GetOnlineDashboard(c.Request().Context())
	if err != nil {
		h.logg.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, response.Data{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, response.Data{
		Message: "Success",
		Data:    sessions,
	})
}

func (h *handlers) GetUserActivity(c echo.Context) error {
	var req request.UserActivity

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

	activity, err := h.svc.GetUserActivity(c.Request().Context(), &req)
	if err != nil {
		h.logg.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, response.Data{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, response.Data{
		Message: "Success",
		Data:    activity,
	})
}
