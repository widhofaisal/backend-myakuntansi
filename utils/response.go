package utils

import (
	"backend-file-management/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func SendSuccess(c echo.Context, message string, data any) error {
	resp := model.ResponseSuccess{
		Success: true,
		Message: message,
		Data:    data,
	}
	return c.JSON(http.StatusOK, resp)
}

func SendError(c echo.Context, status int, message string, err string) error {
	resp := model.ResponseError{
		Success: false,
		Message: message,
		Error:   err,
	}
	return c.JSON(status, resp)
}
