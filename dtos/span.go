package dtos

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type (
	ErrorDTO struct {
		Name    string `json:"name" validate:"required"`
		Message string `json:"message" validate:"required"`
		Stack   string `json:"stack" validate:"required"`
		Cause   string `json:"cause"`
	}
	HTTPDTO struct {
		Resource string `json:"resource" validate:"required"`
		Method   string `json:"method" validate:"required"`
		Headers  string `json:"headers" validate:"required"`
		Status   int    `json:"status"  validate:"required"`
		Data     string `json:"data"`
		Response string `json:"response"`
	}
	EventDTO struct {
		Name  string `json:"name" validate:"required"`
		Value string `json:"value"`
	}
	PerformanceDTO struct {
		Name  string  `json:"name" validate:"required"`
		Value float64 `json:"value" validate:"required"`
		URL   string  `json:"url" validate:"required"`
	}
	PageViewDTO struct {
		URL string `json:"url" validate:"required"`
	}

	BaseSessionDTO struct {
		AppVersion string `json:"appVersion"`
		URL        string `json:"url" validate:"required"`
		Screen     string `json:"screen" validate:"required"`
		Referrer   string `json:"referrer" validate:"required"`
		Language   string `json:"language" validate:"required"`
		Meta       string `json:"meta"`
	}
	SessionDTO struct {
		BaseSessionDTO
		UA string `json:"ua" validate:"required"`
		IP string `json:"ip" validate:"required"`
	}

	CustomValidator struct {
		Validator *validator.Validate
	}
)

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
