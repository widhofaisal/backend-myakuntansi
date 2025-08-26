package controller

import (
	"backend-file-management/utils"

	"github.com/labstack/echo/v4"
)

// Endpoint 1 : Hello
func Hello(c echo.Context) error {
	return utils.SendSuccess(c, "Hello Web File Management by AKUNTANSI BPKAD 2025", 1.3)
}
