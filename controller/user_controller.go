package controller

import (
	"backend-file-management/config"
	"backend-file-management/model"
	"backend-file-management/utils"
	"log"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

// PEGAWAI PAGE

// Endpoint 1 and 2: Get_all_admins_and_employees
func Get_all_admins_and_users(c echo.Context) error {
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	role, _ := utils.Get_role_from_token(token)

	if role != "admin" {
		log.Print(color.RedString("only admins are allowed"))
		return c.JSON((http.StatusUnauthorized), map[string]interface{}{
			"status":  401,
			"message": "unauthorized, only admins are allowed",
		})
	}

	users := []model.User{}
	if err_find := config.DB.Select("id, fullname,username,role").Find(&users).Error; err_find != nil {
		log.Print(color.RedString(err_find.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error, failed to get data",
		})
	}

	type result struct {
		Id       uint   `json:"id"`
		Fullname string `json:"fullname"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	results := []result{}

	for _, user := range users {
		temp := result{
			Id:       user.ID,
			Fullname: user.Fullname,
			Username: user.Username,
			Role:     user.Role,
		}
		results = append(results, temp)
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  200,
		"message": "success to get all data",
		"data":    results,
	})
}

// Endpoint 3 and 4 : Add_admin_and_employee
func Add_admin_and_user(c echo.Context) error {
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	role, _ := utils.Get_role_from_token(token)

	if role != "admin" {
		log.Print(color.RedString("only admins are allowed"))
		return c.JSON((http.StatusUnauthorized), map[string]interface{}{
			"status":  401,
			"message": "unauthorized, only admins are allowed",
		})
	}

	// define struct only for binding
	var bindingUser struct {
		Fullname string `json:"fullname" validate:"required"`
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
		Role     string `json:"role" validate:"required"`
	}

	// binding struct
	if err_bind := c.Bind(&bindingUser); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// check : check: role must be between admin or employee
	accountRole := bindingUser.Role
	if accountRole != "admin" && accountRole != "user" {
		log.Print(color.RedString("bad request, only can add account with role 'admin' or 'user'"))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request, only can add account with role 'admin' or 'user'",
		})
	}

	// check if username already used
	var user model.User
	if err_first := config.DB.Where("username=?", bindingUser.Username).First(&user).Error; err_first == nil {
		log.Print(color.RedString("username already used"))
		return c.JSON((http.StatusConflict), map[string]interface{}{
			"status":  409,
			"message": "status conflic, username already used",
		})
	}

	// check if request body empty
	if bindingUser.Fullname == "" || bindingUser.Username == "" || bindingUser.Password == "" {
		log.Print(color.RedString("request body not valid"))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// hashing password
	hashedPassword, _ := utils.HashPassword(bindingUser.Password)

	// insert data
	user.Fullname = bindingUser.Fullname
	user.Username = bindingUser.Username
	user.Password = hashedPassword
	user.Role = bindingUser.Role
	if err_save := config.DB.Save(&user).Error; err_save != nil {
		log.Print(color.RedString(err_save.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  201,
		"message": "success to register account " + accountRole,
		"data": map[string]interface{}{
			"id":       user.ID,
			"fullname": user.Fullname,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Endpoint 5 and 6 : Update_admin_and_employee
func Update_admin_and_employee(c echo.Context) error {
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	role, _ := utils.Get_role_from_token(token)

	if role != "admin" {
		log.Print(color.RedString("only admins are allowed"))
		return c.JSON((http.StatusUnauthorized), map[string]interface{}{
			"status":  401,
			"message": "unauthorized, only admins are allowed",
		})
	}

	user_id := c.Param("user_id")
	var user model.User

	// check : check: role must be between admin or employee
	accountRole := c.Param("role")
	if accountRole != "admins" && accountRole != "employees" {
		log.Print(color.RedString("bad request, only can update account with role 'admins' or 'employees'"))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request, only can update account with role 'admins' or 'employees'",
		})
	}

	// check is user_id exists and match with accountRole
	if err_first := config.DB.Where("id=? AND role=?", user_id, accountRole[:len(accountRole)-1]).First(&user).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request : there is no user_id = " + user_id + " with role = " + accountRole[:len(accountRole)-1],
		})
	}

	// define struct only for binding
	var bindingUser struct {
		Fullname string `json:"fullname"`
		Password string `json:"password"`
	}

	// binding
	if err_bind := c.Bind(&bindingUser); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request",
		})
	}

	// check if request body empty
	if bindingUser.Fullname == "" {
		log.Print(color.RedString("request body not valid"))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request : request body not valid",
		})
	}

	if bindingUser.Password != "" {
		// hashing password
		hashedPassword, _ := utils.HashPassword(bindingUser.Password)
		user.Password = hashedPassword
	}

	// save
	user.Fullname = bindingUser.Fullname
	// user.Username = bindingUser.Username
	if err_save := config.DB.Save(&user).Error; err_save != nil {
		log.Print(color.RedString(err_save.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"status":  201,
		"message": "success to update " + accountRole[:len(accountRole)-1],
		"data": map[string]interface{}{
			"id":       user.ID,
			"fullname": user.Fullname,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

// Endpoint 7 and 8 : Delete_admin_and_employee
func Delete_admin_and_employee(c echo.Context) error {
	token := strings.TrimPrefix(c.Request().Header.Get("Authorization"), "Bearer ")
	role, _ := utils.Get_role_from_token(token)

	if role != "admin" {
		log.Print(color.RedString("only admins are allowed"))
		return c.JSON((http.StatusUnauthorized), map[string]interface{}{
			"status":  401,
			"message": "unauthorized, only admins are allowed",
		})
	}

	user_id := c.Param("user_id")
	var user model.User

	// check : check: role must be between admin or employee
	accountRole := c.Param("role")
	if accountRole != "admins" && accountRole != "employees" {
		log.Print(color.RedString("bad request, only can update account with role 'admins' or 'employees'"))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request, only can update account with role 'admins' or 'employees'",
		})
	}

	// check is user_id exists and match with accountRole
	if err_first := config.DB.Where("id=? AND role=?", user_id, accountRole[:len(accountRole)-1]).First(&user).Error; err_first != nil {
		log.Print(color.RedString(err_first.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request : there is no user_id = " + user_id + " with role = " + accountRole[:len(accountRole)-1],
		})
	}

	// delete
	if err_delete := config.DB.Where("id=?", user_id).Delete(&user).Error; err_delete != nil {
		log.Print(color.RedString(err_delete.Error()))
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  500,
			"message": "internal server error",
		})
	}

	return c.JSON(http.StatusAccepted, map[string]interface{}{
		"status":  202,
		"message": "success to delete account " + accountRole[:len(accountRole)-1],
		"data": map[string]interface{}{
			"id":       user.ID,
			"fullname": user.Fullname,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}
