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

// CreateProject handles POST /api/auth/project
func CreateProject(c echo.Context) error {
	var input model.Project

	if err := c.Bind(&input); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Bad request", err.Error())
	}

	if input.Name == "" || input.Description == "" {
		log.Print(color.RedString("empty fields in create project"))
		return utils.SendError(c, http.StatusBadRequest, "Bad request", "All fields are required")
	}

	if err := config.DB.Create(&input).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to create project")
	}

	return utils.SendSuccess(c, "Project created successfully", input)
}

// GetAllProjects handles GET /api/auth/project
func GetAllProjects(c echo.Context) error {
	var projects []model.Project

	if err := config.DB.Find(&projects).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to fetch projects")
	}

	return utils.SendSuccess(c, "List of projects", projects)
}

// GetProjectByID handles GET /api/auth/project/:id
func GetProjectByID(c echo.Context) error {
	id := c.Param("id")
	var project model.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusNotFound, "Not found", "Project not found")
	}

	return utils.SendSuccess(c, "Project found", project)
}

// UpdateProject handles PUT /api/auth/project/:id
func UpdateProject(c echo.Context) error {
	id := c.Param("id")
	var project model.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusNotFound, "Not found", "Project not found")
	}

	var input model.Project
	if err := c.Bind(&input); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Bad request", err.Error())
	}

	project.Name = input.Name
	project.Description = input.Description

	if err := config.DB.Save(&project).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to update project")
	}

	return utils.SendSuccess(c, "Project updated successfully", project)
}

// DeleteProject handles DELETE /api/auth/project/:id
func DeleteProject(c echo.Context) error {
	id := c.Param("id")
	var project model.Project

	if err := config.DB.First(&project, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusNotFound, "Not found", "Project not found")
	}

	if err := config.DB.Delete(&project).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to delete project")
	}

	return utils.SendSuccess(c, "Project deleted successfully", nil)
}
