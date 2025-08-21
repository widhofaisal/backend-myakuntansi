package controller

import (
	"backend-file-management/config"
	"backend-file-management/model"
	"backend-file-management/utils"
	"fmt"
	"log"
	"net/http"

	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
)

// CreateItem handles POST /api/auth/item
func CreateItem(c echo.Context) error {
	var input model.Item

	if err := c.Bind(&input); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Bad request", err.Error())
	}

	// Validasi field wajib
	if input.Name == "" || input.Type == "" {
		log.Print(color.RedString("empty fields in create item"))
		return utils.SendError(c, http.StatusBadRequest, "Bad request", "Name and Type are required")
	}

	// Validasi unique
	if err := config.DB.Where("name = ?", input.Name).Where("parent_id = ?", input.ParentID).First(&model.Item{}).Error; err == nil {
		log.Print(color.RedString("item name already exists"))
		return utils.SendError(c, http.StatusBadRequest, "Bad request", "Item name already exists")
	}

	// Validasi tipe yang diperbolehkan
	allowedTypes := map[model.ItemType]bool{
		model.ItemTypeFolder: true,
		model.ItemTypePdf:    true,
		model.ItemTypeJpg:    true,
		model.ItemTypePng:    true,
	}
	if !allowedTypes[input.Type] {
		log.Print(color.RedString("invalid item type"))
		return utils.SendError(c, http.StatusBadRequest, "Bad request", "Invalid item type")
	}

	// Validasi parent_id (jika tidak null)
	if input.ParentID != nil {
		var parent model.Item
		if err := config.DB.First(&parent, *input.ParentID).Error; err != nil {
			log.Print(color.RedString("invalid parent_id"))
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Parent folder not found")
		}
		// Pastikan parent memang folder
		if parent.Type != model.ItemTypeFolder {
			log.Print(color.RedString("parent is not a folder"))
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Parent must be a folder")
		}
	}

	// Jika type folder, maka FilePath, MimeType, Size harus null
	if input.Type == model.ItemTypeFolder {
		input.FilePath = nil
		input.MimeType = nil
		input.Size = nil
	}

	if err := config.DB.Create(&input).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to create item")
	}

	return utils.SendSuccess(c, "Item created successfully", input)
}

// GetAllItems handles GET /api/auth/item
func GetAllItems(c echo.Context) error {
	var items []model.Item

	if err := config.DB.Preload("Parent").Find(&items).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to fetch items")
	}

	return utils.SendSuccess(c, "List of items", items)
}

// GetAllItemsAndFolders handles GET /api/auth/item
func GetAllItemsAndFolders(c echo.Context) error {
	id := c.Param("id")
	var folders []model.Item
	var files []model.Item
	var currentPath string
	var breadcrumbs []map[string]any

	if id == "0" {
		if err := config.DB.Where("parent_id IS NULL").Where("is_folder = ?", true).Find(&folders).Error; err != nil {
			log.Print(color.RedString(err.Error()))
			return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to fetch items")
		}
		if err := config.DB.Where("parent_id IS NULL").Where("is_folder = ?", false).Find(&files).Error; err != nil {
			log.Print(color.RedString(err.Error()))
			return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to fetch items")
		}
	} else {
		if err := config.DB.Where("parent_id = ?", id).Where("is_folder = ?", true).Find(&folders).Error; err != nil {
			log.Print(color.RedString(err.Error()))
			return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to fetch items")
		}
		if err := config.DB.Where("parent_id = ?", id).Where("is_folder = ?", false).Find(&files).Error; err != nil {
			log.Print(color.RedString(err.Error()))
			return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to fetch items")
		}
	}

	rawFindCurrentPath := `	WITH RECURSIVE path_cte AS (
								SELECT 
									id,
									parent_id,
									name,
									CAST(name AS CHAR(1000)) AS full_path
								FROM items
								WHERE id = ?

								UNION ALL

								SELECT 
									t.id,
									t.parent_id,
									t.name,
									CONCAT(t.name, '/', c.full_path) AS full_path
								FROM items t
								INNER JOIN path_cte c ON t.id = c.parent_id
							)
							SELECT CONCAT('/', full_path) AS path
							FROM path_cte
							WHERE parent_id IS NULL;
						`
	if err := config.DB.Raw(rawFindCurrentPath, id).Scan(&currentPath).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to fetch items")
	}

	rawFindBreadCrumbs := `	WITH RECURSIVE cte (id, name, parent_id) AS (
								SELECT id, name, parent_id
								FROM items
								WHERE id = ?
								UNION ALL
								SELECT t.id, t.name, t.parent_id
								FROM items t
								JOIN cte ON t.id = cte.parent_id
							)
							SELECT 
								id,
								name
							FROM cte
							ORDER BY id;
							`
	if err := config.DB.Raw(rawFindBreadCrumbs, id).Scan(&breadcrumbs).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to fetch items")
	}

	return utils.SendSuccess(c, "List of items and folders", map[string]any{
		"folders":     folders,
		"files":       files,
		"currentPath": currentPath,
		"breadcrumbs": breadcrumbs,
	})
}

// GetItemByID handles GET /api/auth/item/:id
func GetItemByID(c echo.Context) error {
	id := c.Param("id")
	var item model.Item

	if err := config.DB.Preload("Parent").First(&item, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusNotFound, "Not found", "Item not found")
	}

	return utils.SendSuccess(c, "Item found", item)
}

// UpdateItem handles PUT /api/auth/item/:id
func UpdateItem(c echo.Context) error {
	id := c.Param("id")
	var item model.Item

	if err := config.DB.First(&item, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusNotFound, "Not found", "Item not found")
	}

	var input model.Item
	if err := c.Bind(&input); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Bad request", err.Error())
	}

	// Update field yang diizinkan
	if input.Name != "" {
		item.Name = input.Name
	}
	if input.Type != "" {
		item.Type = input.Type
		if input.Type == model.ItemTypeFolder {
			item.FilePath = nil
			item.MimeType = nil
			item.Size = nil
		} else {
			item.FilePath = input.FilePath
			item.MimeType = input.MimeType
			item.Size = input.Size
		}
	}
	if input.ParentID != nil {
		item.ParentID = input.ParentID
	}
	if input.UploadedBy != nil {
		item.UploadedBy = input.UploadedBy
	}

	if err := config.DB.Save(&item).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to update item")
	}

	return utils.SendSuccess(c, "Item updated successfully", item)
}

// DeleteItem handles DELETE /api/auth/item/:id
func DeleteItem(c echo.Context) error {
	id := c.Param("id")
	var item model.Item

	if err := config.DB.First(&item, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusNotFound, "Not found", "Item not found")
	}

	if err := config.DB.Delete(&item).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to delete item")
	}

	return utils.SendSuccess(c, "Item deleted successfully", nil)
}

// DownloadFile handles GET /api/auth/items/download/:id
func DownloadFile(c echo.Context) error {
	id := c.Param("id")
	fmt.Println(id)
	var item model.Item

	if err := config.DB.First(&item, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusNotFound, "Not found", "Item not found")
	}
	fmt.Println("masuk1")

	if item.Type != model.ItemTypePdf && item.Type != model.ItemTypeJpg && item.Type != model.ItemTypePng {
		log.Print(color.RedString("invalid file type"))
		return utils.SendError(c, http.StatusBadRequest, "Bad request", "Invalid file type")
	}
	fmt.Println(*item.FilePath)

	return c.File(*item.FilePath)
}
