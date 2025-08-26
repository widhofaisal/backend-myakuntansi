package controller

import (
	"backend-file-management/config"
	"backend-file-management/model"
	"backend-file-management/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// CreateItem handles POST /api/auth/item
// func CreateItem(c echo.Context) error {
// 	var input model.Item

// 	// Ambil field dari form
// 	input.Name = c.FormValue("name")
// 	input.Type = model.ItemType(c.FormValue("type"))

// 	if parentID := c.FormValue("parent_id"); parentID != "" {
// 		id, err := strconv.Atoi(parentID)
// 		if err != nil {
// 			log.Println("Invalid parent_id:", err)
// 			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Invalid parent_id")
// 		}
// 		uid := uint(id)
// 		input.ParentID = &uid
// 	}

// 	// Debug
// 	log.Println("=== Form Values ===")
// 	log.Println("Name:", input.Name)
// 	log.Println("Type:", input.Type)
// 	log.Println("ParentID:", input.ParentID)

// 	// Validasi field wajib
// 	if input.Name == "" || input.Type == "" {
// 		return utils.SendError(c, http.StatusBadRequest, "Bad request", "Name and Type are required")
// 	}

// 	// Kalau bukan folder → wajib ada file
// 	if input.Type != model.ItemTypeFolder {
// 		fileHeader, err := c.FormFile("file")
// 		if err != nil {
// 			log.Println("No file uploaded:", err)
// 			return utils.SendError(c, http.StatusBadRequest, "Bad request", "File is required for non-folder items")
// 		}

// 		// Debug
// 		log.Printf("Uploaded file: %s (Header: %+v)\n", fileHeader.Filename, fileHeader.Header)

// 		// Validasi ukuran file (maks 10 MB)
// 		if fileHeader.Size > 10<<20 {
// 			return utils.SendError(c, http.StatusBadRequest, "Bad request", "File size exceeds 10MB limit")
// 		}

// 		// Buat folder upload kalau belum ada
// 		uploadDir := "uploads"
// 		if err := os.MkdirAll(uploadDir, 0755); err != nil {
// 			log.Println("Failed to create upload dir:", err)
// 			return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to create uploads directory")
// 		}

// 		// Simpan file dengan UUID
// 		extension := filepath.Ext(fileHeader.Filename)
// 		newFilename := uuid.New().String() + extension
// 		destination := filepath.Join(uploadDir, newFilename)

// 		log.Println("Saving file to:", destination)

// 		if err := c.SaveUploadedFile(fileHeader, destination); err != nil {
// 			log.Println("SaveUploadedFile error:", err)
// 			return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to save file")
// 		}

// 		// Isi metadata file
// 		mimeType := fileHeader.Header.Get("Content-Type")
// 		size := fileHeader.Size

// 		log.Println("File saved successfully")
// 		log.Println("MimeType:", mimeType)
// 		log.Println("Size:", size)

// 		input.FilePath = &destination
// 		input.MimeType = &mimeType
// 		input.Size = &size
// 	} else {
// 		input.FilePath = nil
// 		input.MimeType = nil
// 		input.Size = nil
// 	}

// 	// Validasi nama unik dalam parent folder
// 	if err := config.DB.Where("name = ?", input.Name).Where("parent_id = ?", input.ParentID).First(&model.Item{}).Error; err == nil {
// 		return utils.SendError(c, http.StatusBadRequest, "Bad request", "Item name already exists in this location")
// 	}

// 	// Validasi tipe item
// 	allowedTypes := map[model.ItemType]bool{
// 		model.ItemTypeFolder: true,
// 		model.ItemTypePdf:    true,
// 		model.ItemTypeJpg:    true,
// 		model.ItemTypePng:    true,
// 	}
// 	if !allowedTypes[input.Type] {
// 		return utils.SendError(c, http.StatusBadRequest, "Bad request", "Invalid item type")
// 	}

// 	// Validasi parent
// 	if input.ParentID != nil {
// 		var parent model.Item
// 		if err := config.DB.First(&parent, *input.ParentID).Error; err != nil {
// 			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Parent folder not found")
// 		}
// 		if parent.Type != model.ItemTypeFolder {
// 			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Parent must be a folder")
// 		}
// 	}

// 	// Simpan ke database
// 	if err := config.DB.Create(&input).Error; err != nil {
// 		// Kalau gagal, hapus file yang sudah diupload
// 		if input.FilePath != nil {
// 			_ = os.Remove(*input.FilePath)
// 		}
// 		log.Println("DB Create error:", err)
// 		return utils.SendError(c, http.StatusInternalServerError, "Database error", "Failed to create item")
// 	}

// 	log.Println("Item created successfully:", input)

// 	return utils.SendSuccess(c, "Item created successfully", input)
// }

func CreateFile(c echo.Context) error {
	fmt.Println("CreateFile")
	// Ambil form value
	name := c.FormValue("name")
	parentIDStr := c.FormValue("parentId")

	// Parse parent_id (boleh kosong/null)
	var parentID *uint
	if parentIDStr != "" {
		idParsed, err := strconv.ParseUint(parentIDStr, 10, 64)
		if err != nil {
			return utils.SendError(c, http.StatusBadRequest, "Bad request", "Invalid parent_id")
		}
		idUint := uint(idParsed)
		parentID = &idUint
	}

	// Ambil file dari form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		log.Println("No file uploaded:", err)
		return utils.SendError(c, http.StatusBadRequest, "Bad request", "File is required")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to open file")
	}
	defer src.Close()

	// Buat nama file dengan UUID + ekstensi asli
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		ext = ".bin"
	}
	fileName := uuid.NewString() + ext
	filePath := filepath.Join("uploads", fileName)

	// Pastikan folder uploads ada
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to create uploads folder")
	}

	// Simpan file
	dst, err := os.Create(filePath)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to create file")
	}
	defer dst.Close()

	size, err := io.Copy(dst, src)
	if err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to save file")
	}

	// Ambil mime type
	mimeType := fileHeader.Header.Get("Content-Type")

	// Tentukan tipe berdasarkan ekstensi
	var itemType model.ItemType
	switch ext {
	case ".pdf":
		itemType = model.ItemTypePdf
	case ".jpg", ".jpeg":
		itemType = model.ItemTypeJpg
	case ".png":
		itemType = model.ItemTypePng
	default:
		itemType = model.ItemTypeOther // semua jenis file lain masuk sini
	}

	// Simpan metadata ke DB
	item := model.Item{
		Name:     name,
		ParentID: parentID,
		Type:     itemType,
		FilePath: &fileName,
		MimeType: &mimeType,
		Size:     &size,
		IsFolder: false,
		IsLink:   false,
	}
	if err := config.DB.Create(&item).Error; err != nil {
		return utils.SendError(c, http.StatusInternalServerError, "Server error", "Failed to save item to DB")
	}

	return utils.SendSuccess(c, "File created successfully", item)
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

func CreateFolder(c echo.Context) error {
	var folderRequest struct {
		name       string `json:"name"`
		parentId   *int   `json:"parentId"`
		uploadedBy string `json:"uploadedBy"`
	}

	// binding request body
	if err_bind := c.Bind(&folderRequest); err_bind != nil {
		log.Print(color.RedString(err_bind.Error()))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request, request body not valid",
		})
	}

	// validate empty
	if folderRequest.name == "" || folderRequest.parentId == nil || folderRequest.uploadedBy == "" {
		log.Print(color.RedString("name, parentId, uploadedBy cannot be empty"))
		return c.JSON((http.StatusBadRequest), map[string]interface{}{
			"status":  400,
			"message": "bad request, name, parentId, uploadedBy cannot be empty",
		})
	}

	// insert data
	parentID := uint(*folderRequest.parentId)
	item := model.Item{
		Name:       folderRequest.name,
		ParentID:   &parentID,
		Type:       model.ItemTypeFolder,
		UploadedBy: nil,
		IsFolder:   true,
		IsLink:     false,
	}

	if err := config.DB.Create(&item).Error; err != nil {
		log.Print(color.RedString(err.Error()))
		return c.JSON((http.StatusInternalServerError), map[string]interface{}{
			"status":  500,
			"message": "internal server error, failed to create folder",
		})
	}

	return c.JSON((http.StatusOK), map[string]interface{}{
		"status":  200,
		"message": "success create folder",
	})
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
