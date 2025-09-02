package controller

import (
	"backend-file-management/config"
	"backend-file-management/model"
	"backend-file-management/utils"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

func getUploadsSize() (int64, error) {
	var totalSize int64
	err := filepath.Walk("uploads", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			totalSize += info.Size()
		}
		return nil
	})
	return totalSize, err
}

func CountDashboard(c echo.Context) error {
	var countDashboard struct {
		TotalFolders     int64            `json:"totalFolders"`
		TotalFiles       int64            `json:"totalFiles"`
		TotalLinks       int64            `json:"totalLinks"`
		TotalStorageUsed int64            `json:"totalStorageUsed"`
		TotalUsers       int64            `json:"totalUsers"`
		FilesByExtension map[string]int64 `json:"filesByExtension"`
	}

	if err := config.DB.Model(&model.Item{}).Where("is_folder = ?", true).Count(&countDashboard.TotalFolders).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to count folders")
	}

	if err := config.DB.Model(&model.Item{}).Where("is_folder = ? AND is_link = ?", false, false).Count(&countDashboard.TotalFiles).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to count files")
	}

	if err := config.DB.Model(&model.Item{}).Where("is_link = ?", true).Count(&countDashboard.TotalLinks).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to count links")
	}

	if err := config.DB.Model(&model.User{}).Count(&countDashboard.TotalUsers).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to count users")
	}

	// Calculate total storage used in uploads directory
	totalSize, err := getUploadsSize()
	if err != nil {
		log.Print(color.RedString("Error calculating uploads size: " + err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to calculate storage usage")
	}
	countDashboard.TotalStorageUsed = totalSize

	// Count files by extension
	var filesByExtension []struct {
		MimeType string `gorm:"column:mime_type"`
		Count    int64  `gorm:"column:count"`
	}

	if err := config.DB.Model(&model.Item{}).
		Select("mime_type, COUNT(*) as count").
		Where("is_folder = ? AND is_link = ?", false, false).
		Group("mime_type").
		Find(&filesByExtension).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to count files by extension")
	}

	// Convert the result to a map
	countDashboard.FilesByExtension = make(map[string]int64)
	for _, item := range filesByExtension {
		countDashboard.FilesByExtension[item.MimeType] = item.Count
	}

	return utils.SendSuccess(c, "Dashboard data", countDashboard)
}
