package route

import (
	"backend-file-management/constant"
	"backend-file-management/controller"
	"backend-file-management/middleware"

	"github.com/labstack/echo/v4"
	mid "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	// Middleware global
	e.Use(mid.CORS())
	e.Use(middleware.MiddlewareLogging)

	// Root API group: /api/v1
	api := e.Group("/api/v1")

	// ================================
	// --------- PUBLIC ROUTES -------
	// ================================
	api.GET("/hello", controller.Hello)
	api.POST("/login", controller.Login)	

	// ================================
	// --------- AUTH ROUTES ---------
	// ================================
	auth := api.Group("/auth")
	auth.Use(JWTMiddleware())

	auth.GET("/users", controller.Get_all_admins_and_users)
	auth.POST("/users", controller.Add_admin_and_user)
	
	auth.POST("/projects", controller.CreateProject)
	auth.GET("/projects", controller.GetAllProjects)
	auth.GET("/projects/:id", controller.GetProjectByID)
	auth.PUT("/projects/:id", controller.UpdateProject)
	auth.DELETE("/projects/:id", controller.DeleteProject)

	// auth.POST("/items", controller.CreateItem)
	auth.POST("/file", controller.CreateFile)
	auth.POST("/folder", controller.CreateFolder)
	// auth.GET("/items_and_folders/:id", controller.GetAllItemsAndFolders)
	auth.GET("/items", controller.GetAllItems)
	auth.GET("/items/:id", controller.GetItemByID)
	auth.PUT("/items/:id", controller.UpdateItem)
	auth.DELETE("/items/:id", controller.DeleteItem)
	auth.GET("/items/download/:id", controller.DownloadFile)

	// Contoh: auth.GET("/shift", controller.Get_all_shift)

	return e
}

func JWTMiddleware() echo.MiddlewareFunc {
	config := mid.JWTConfig{
		SigningKey: []byte(constant.SECRET_JWT),
	}
	return mid.JWTWithConfig(config)
}
