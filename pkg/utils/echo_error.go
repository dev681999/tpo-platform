package utils

import (
	"github.com/labstack/echo/v4"
)

// NewEchoErrorResponse is a helper function to return http error
func NewEchoErrorResponse(statusCode int, err error) error {
	return echo.NewHTTPError(statusCode, err.Error())
}
