package v1

import (
	"errors"

	"github.com/labstack/echo/v4"
)

var (
	ErrInternalServer        = errors.New("internal server error")
	ErrInvalidRequestBody    = errors.New("invalid request body")
	ErrCommandNotFound       = errors.New("command not found")
	ErrInvalidPathParameter  = errors.New("invalid path parameter")
	ErrCommandAlreadyRunning = errors.New("command already running")
	ErrCommandNotRunning     = errors.New("command not currenlty running")
)

func newErrorMessage(c echo.Context, statusCode int, message string) error {
	httpErr := echo.NewHTTPError(statusCode, message)
	return c.JSON(statusCode, httpErr)
}
