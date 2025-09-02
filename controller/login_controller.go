package controller

import (
	"backend-file-management/config"
	"backend-file-management/model"
	"backend-file-management/utils"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	var user model.User

	// define struct only for binding
	var bindingUser struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	// binding struct
	if err := c.Bind(&bindingUser); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Bad request", err.Error())
	}

	// check if request body empty
	if bindingUser.Username == "" || bindingUser.Password == "" {
		log.Print(color.RedString("request body not valid"))
		return utils.SendError(c, http.StatusBadRequest, "Bad request", "Username and password are required")
	}

	// check username validity
	if err := config.DB.Where("username = ?", bindingUser.Username).First(&user).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", "Username not found")
	}

	// verify the password
	if !utils.CheckPasswordHash(bindingUser.Password, user.Password) {
		log.Print(color.RedString("password is incorrect"))
		return utils.SendError(c, http.StatusUnauthorized, "Unauthorized", "Invalid password")
	}

	// create token
	token, err := utils.Create_token(uint(user.ID), user.Username, user.Role)
	if err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Internal server error", "Failed to generate token")
	}

	// response success
	data := map[string]any{
		"token": token,
		"user": model.User{
			ID:       user.ID,
			Username: user.Username,
			Role:     user.Role,
		},
	}
	return utils.SendSuccess(c, "Success to login", data)
}
