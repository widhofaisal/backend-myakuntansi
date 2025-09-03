package controller

import (
	"backend-file-management/config"
	"backend-file-management/model"
	"backend-file-management/utils"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateLink(c echo.Context) error {
	var folderRequest struct {
		Name     string `json:"name"`
		ParentId *int   `json:"parentId,omitempty"`
		FilePath string `json:"filePath"`
	}

	// binding request body
	if err_bind := c.Bind(&folderRequest); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request, request body not valid",
		})
	}
	fmt.Println("folderRequest:", folderRequest)

	// validate name
	if folderRequest.Name == "" || folderRequest.FilePath == "" {
		log.Print(color.RedString("name and filepath cannot be empty"))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request, name and filepath cannot be empty",
		})
	}

	// create item with optional parent
	item := model.Item{
		Name:     folderRequest.Name,
		Type:     model.ItemTypeFolder,
		IsFolder: false,
		IsLink:   true,
		FilePath: &folderRequest.FilePath,
	}

	parentIDStr := folderRequest.ParentId
	if parentIDStr != nil && *parentIDStr != 0 {

		// Validate parent exists and is a folder
		var parent model.Item
		if err := config.DB.First(&parent, *parentIDStr).Error; err != nil {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Parent folder not found2")
		}
		if !parent.IsFolder {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Parent must be a folder")
		}

		// Check for duplicate name in the same folder
		var existingItem model.Item
		if err := config.DB.Where("name = ? AND parent_id = ?", folderRequest.Name, *parentIDStr).First(&existingItem).Error; err == nil {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "A file with this name already exists in the destination folder")
		}
	} else {
		// Check for duplicate name in root folder (where parent_id is NULL)
		var existingItem model.Item
		if err := config.DB.Where("name = ? AND parent_id IS NULL", folderRequest.Name).First(&existingItem).Error; err == nil {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "A file with this name already exists in the root folder")
		}
	}

	// Save to database
	if err := config.DB.Create(&item).Error; err != nil {
		log.Printf("Error creating link: %v", err)
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to create link")
	}

	return utils.SendSuccess(c, "link created successfully", item)
}

func CreateFile(c echo.Context) error {
	// Get file from form first
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Println("No file uploaded:", err)
		return utils.SendError(c, http.StatusBadRequest, "Bad request", "File is required")
	}

	// Get parentId from form (can be empty/null)
	parentIDStr := c.FormValue("parentId")
	originalName := filepath.Base(fileHeader.Filename)

	// Parse parent_id (can be empty/null)
	var parentID *uint
	if parentIDStr != "" && parentIDStr != "0" {
		idParsed, err := strconv.ParseUint(parentIDStr, 10, 64)
		if err != nil {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Invalid parent_id")
		}
		idUint := uint(idParsed)
		parentID = &idUint

		// Validate parent exists and is a folder
		var parent model.Item
		if err := config.DB.First(&parent, parentID).Error; err != nil {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Parent folder not found")
		}
		if !parent.IsFolder {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Parent must be a folder")
		}

		// Check for duplicate name in the same folder
		var existingItem model.Item
		if err := config.DB.Where("name = ? AND parent_id = ?", originalName, parentID).First(&existingItem).Error; err == nil {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "A file with this name already exists in the destination folder")
		}
	} else {
		// Check for duplicate name in root folder (where parent_id is NULL)
		var existingItem model.Item
		if err := config.DB.Where("name = ? AND parent_id IS NULL", originalName).First(&existingItem).Error; err == nil {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "A file with this name already exists in the root folder")
		}
	}

	src, err := fileHeader.Open()
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to open file")
	}
	defer src.Close()

	// Create unique filename with UUID + original extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".bin"
	}
	fileName := uuid.NewString() + ext
	filePath := filepath.Join("uploads", fileName)

	// Ensure uploads directory exists
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to create uploads folder")
	}

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to create file")
	}
	defer dst.Close()

	size, err := io.Copy(dst, src)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to save file")
	}

	// Get mime type
	mimeType := fileHeader.Header.Get("Content-Type")

	// Determine item type based on extension
	var itemType model.ItemType
	switch ext {
	case ".pdf":
		itemType = model.ItemTypePdf
	case ".jpg", ".jpeg":
		itemType = model.ItemTypeJpg
	case ".png":
		itemType = model.ItemTypePng
	default:
		itemType = model.ItemTypeOther
	}

	// Save metadata to DB
	item := model.Item{
		Name:     originalName, // Use original filename as the name
		ParentID: parentID,
		Type:     itemType,
		FilePath: &fileName,
		MimeType: &mimeType,
		Size:     &size,
		IsFolder: false,
		IsLink:   false,
	}

	if err := config.DB.Create(&item).Error; err != nil {
		// Clean up the uploaded file if DB save fails
		os.Remove(filePath)
		return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to save item to DB")
	}

	return utils.SendSuccess(c, "File uploaded successfully", item)
}

func GetAllItems(c echo.Context) error {
	var items []model.Item

	if err := config.DB.Preload("Parent").Find(&items).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to fetch items")
	}

	return utils.SendSuccess(c, "List of items", items)
}

func CreateFolder(c echo.Context) error {
	var folderRequest struct {
		Name     string `json:"name"`
		ParentId *int   `json:"parentId,omitempty"`
	}

	// binding request body
	if err_bind := c.Bind(&folderRequest); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request, request body not valid",
		})
	}

	// validate name
	if folderRequest.Name == "" {
		log.Print(color.RedString("name cannot be empty"))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request, name cannot be empty",
		})
	}

	// create item with optional parent
	item := model.Item{
		Name:     folderRequest.Name,
		Type:     model.ItemTypeFolder,
		IsFolder: true,
		IsLink:   false,
	}

	// Set ParentID only if provided and not nil
	parentIDStr := c.FormValue("parentId")
	if parentIDStr != "" && parentIDStr != "0" {
		fmt.Println("parentId tidak kosong")
		parentID := uint(*folderRequest.ParentId)
		item.ParentID = &parentID

		// Validasi parent
		var parent model.Item
		if err := config.DB.First(&parent, *folderRequest.ParentId).Error; err != nil {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Parent folder not found")
		}
		if parent.Type != model.ItemTypeFolder {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Parent must be a folder")
		}

		// Validasi nama unik dalam parent folder
		if err := config.DB.Where("name = ?", folderRequest.Name).Where("parent_id = ?", folderRequest.ParentId).First(&model.Item{}).Error; err == nil {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Item name already exists in this location")
		}
	} else {
		fmt.Println("parentId kosong")
		// Validasi nama unik dalam parent folder
		if err := config.DB.Where("name = ? AND parent_id IS NULL", folderRequest.Name).First(&model.Item{}).Error; err == nil {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Item name already exists in this location")
		}
		// Explicitly set to nil if not provided
		item.ParentID = nil
	}

	// Save to database
	if err := config.DB.Create(&item).Error; err != nil {
		log.Printf("Error creating folder: %v", err)
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to create folder")
	}

	return utils.SendSuccess(c, "Folder created successfully", item)
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

	var requestJson struct {
		Name string `json:"name"`
	}

	if err := config.DB.First(&item, id).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusNotFound, "Not found", "Item not found")
	}

	if err := c.Bind(&requestJson); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Bad request", err.Error())
	}

	if requestJson.Name != "" {
		item.Name = requestJson.Name
	}

	if err := config.DB.Save(&item).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to update item")
	}

	return utils.SendSuccess(c, "Item updated successfully", item)
}

func deleteItemRecursive(tx *gorm.DB, itemID string) error {
	// Find children of the current item
	var children []model.Item
	if err := tx.Where("parent_id = ?", itemID).Find(&children).Error; err != nil {
		return err
	}

	// Recursively delete each child
	for _, child := range children {
		childIDStr := strconv.FormatUint(uint64(child.ID), 10)
		if err := deleteItemRecursive(tx, childIDStr); err != nil {
			return err
		}
	}

	// After deleting all children, delete the item itself
	if err := tx.Where("id = ?", itemID).Delete(&model.Item{}).Error; err != nil {
		return err
	}

	return nil
}

// DeleteItem handles DELETE /api/auth/item/:id
func DeleteItem(c echo.Context) error {
	id := c.Param("id")

	// Start a transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		log.Print(color.RedString(tx.Error.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to start transaction")
	}

	// Recursively delete the item and its children
	if err := deleteItemRecursive(tx, id); err != nil {
		tx.Rollback()
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to delete item and its children")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		log.Print(color.RedString(err.Error()))
		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to commit transaction")
	}

	return utils.SendSuccess(c, "Item and all its children deleted successfully", nil)
}

// DownloadFile handles GET /api/auth/items/download/:id
func DownloadFile(c echo.Context) error {
	id := c.Param("id")
	log.Printf("Processing download request for file ID: %s", id)

	var item model.Item

	// Get file info from database
	if err := config.DB.First(&item, id).Error; err != nil {
		log.Printf("Error finding file with ID %s: %v", id, err)
		return utils.SendError(c, http.StatusNotFound, "Not found", "File not found")
	}

	if item.FilePath == nil || *item.FilePath == "" {
		errMsg := "File path is empty"
		log.Print(color.RedString(errMsg))
		return utils.SendError(c, http.StatusBadRequest, "Bad request", errMsg)
	}

	// Construct full file path
	filePath := filepath.Clean(*item.FilePath)
	if !strings.HasPrefix(filePath, "uploads") {
		filePath = filepath.Join("uploads", filePath)
	}

	// Get file info to verify it exists
	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		errMsg := fmt.Sprintf("File not found at path: %s", filePath)
		log.Print(color.RedString(errMsg))
		return utils.SendError(c, http.StatusNotFound, "Not found", "File not found on server")
	} else if err != nil {
		errMsg := fmt.Sprintf("Error accessing file: %v", err)
		log.Print(color.RedString(errMsg))
		return utils.SendError(c, http.StatusInternalServerError, "Internal server error", "Error accessing file")
	}

	// Get file name and extension
	fileName := filepath.Base(filePath)
	ext := strings.ToLower(filepath.Ext(fileName))

	// Set default MIME type if extension is empty
	if ext == "" {
		ext = filepath.Ext(fileName) // Try to get extension again without forcing lowercase
	}

	// Set content type
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	// Log the details for debugging
	log.Printf("Serving file - Path: %s, Name: %s, Ext: %s, Type: %s, Size: %d bytes",
		filePath, fileName, ext, contentType, fileInfo.Size())

	// Set headers before sending the file
	c.Response().Header().Set(echo.HeaderContentType, contentType)
	c.Response().Header().Set(echo.HeaderContentDisposition,
		fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	// Send the file
	return c.File(filePath)
}
