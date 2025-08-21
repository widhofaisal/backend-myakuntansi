package utils

import (
	// "fmt"
	// "image"
	// "image/color"
	// "image/draw"
	// "image/jpeg"
	// "backend-file-management/config"
	// "backend-file-management/model"
	// "log"

	// "os"

	// col "github.com/fatih/color"
)

// TODO:
// func AutoInsertShift() bool {
// 	if err_first := config.DB.Find(&model.Shift{}).Error; err_first != nil {
// 		log.Print(col.RedString("CHEAT"), " - already auto insert shift")
// 		return false
// 	}

// 	shifts := [][]interface{}{
// 		{"Shift 1", "07:00 - 15:00"},
// 		{"Shift 2", "15:00 - 23:00"},
// 		{"Shift 3", "23:00 - 07:00"},
// 	}

// 	var shiftModels []model.Shift
// 	for _, shiftData := range shifts {
// 		shiftModels = append(shiftModels, model.Shift{
// 			Name:        shiftData[0].(string),
// 			Description: shiftData[1].(string),
// 		})
// 	}

// 	if err_insert := config.DB.Save(&shiftModels).Error; err_insert != nil {
// 		log.Print(col.RedString("CHEAT"), "- failed to auto insert shift")
// 		return false
// 	}

// 	log.Print(col.RedString("CHEAT"), "- success to auto insert shift")
// 	return true
// }

// TODO:
// func AutoInsertMaterials() bool {
// 	var material model.Material
// 	if err_firstMaterial := config.DB.First(&material).Error; err_firstMaterial == nil {
// 		log.Print(col.RedString("CHEAT"), " - materials already exist")
// 		return false
// 	}

// 	materials := [][]interface{}{
// 		{"Semen", "sak"},
// 		{"Pasir", "m3"},
// 		{"Batu Split", "m3"},
// 		{"Baja Tulangan", "ton"},
// 		{"Bata Merah", "batang"},
// 		{"Besi Beton", "batang"},
// 		{"Kayu Balok", "m3"},
// 		{"Cat Interior", "liter"},
// 		{"Cat Eksterior", "liter"},
// 		{"Genteng", "buah"},
// 		{"Kaca", "m2"},
// 		{"Pintu Kayu", "buah"},
// 		{"Jendela Aluminium", "buah"},
// 		{"Keramik Lantai", "m2"},
// 		{"Keramik Dinding", "m2"},
// 		{"Paku", "kg"},
// 		{"Pipa PVC", "m"},
// 		{"Pipa Besi", "m"},
// 		{"Insulasi", "m2"},
// 		{"Plaster", "m2"},
// 		{"Atap Baja Ringan", "m2"},
// 		{"Rangka Plafon", "m"},
// 		{"Pintu Kaca", "buah"},
// 		{"Cat Pelapis Anti Karat", "liter"},
// 	}

// 	var materialModels []model.Material
// 	for _, materialData := range materials {
// 		materialModels = append(materialModels, model.Material{
// 			MaterialName: materialData[0].(string),
// 			Unit:         materialData[1].(string),
// 		})
// 	}

// 	if errInsert := config.DB.Create(&materialModels).Error; errInsert != nil {
// 		log.Print(col.RedString("CHEAT "), "internal server error, failed to auto insert materials")
// 		return false
// 	} else {
// 		log.Print(col.RedString("CHEAT "), "success auto insert materials")
// 		return true
// 	}

// }

// TODO:
// func AutoCreateFolderReceipts() bool {
// 	// CHECK : is assets/payments exist
// 	dirParent1 := "assets"
// 	dirChild1 := "payments"
// 	fullDir1 := dirParent1 + string(os.PathSeparator) + dirChild1
// 	if _, err := os.Stat(fullDir1); err == nil {
// 		log.Print(col.RedString("CHEAT"), " - directory /assets/payments already exist")
// 	} else if err_create := os.Mkdir(fullDir1, 0755); err_create != nil {
// 		log.Print(col.RedString("CHEAT"), " - error while auto create folder /assets/payments")
// 	}

// 	// CHECK : is assets/invoices exist
// 	dirParent2 := "assets"
// 	dirChild2 := "invoices"
// 	fullDir2 := dirParent2 + string(os.PathSeparator) + dirChild2
// 	if _, err := os.Stat(fullDir2); err == nil {
// 		log.Print(col.RedString("CHEAT"), " - directory /assets/invoices already exist")
// 	} else if err_create := os.Mkdir(fullDir2, 0755); err_create != nil {
// 		log.Print(col.RedString("CHEAT"), " - error while auto create folder /assets/invoices")
// 	}

// 	// CHECK : is assets/taxs exist
// 	dirParent3 := "assets"
// 	dirChild3 := "taxs"
// 	fullDir3 := dirParent3 + string(os.PathSeparator) + dirChild3
// 	if _, err := os.Stat(fullDir3); err == nil {
// 		log.Print(col.RedString("CHEAT"), " - directory /assets/taxs already exist")
// 	} else if err_create := os.Mkdir(fullDir3, 0755); err_create != nil {
// 		log.Print(col.RedString("CHEAT"), " - error while auto create folder /assets/taxs")
// 	}

// 	// CHECK : is assets/delivery_receipts exist
// 	dirParent4 := "assets"
// 	dirChild4 := "delivery_receipts"
// 	fullDir4 := dirParent4 + string(os.PathSeparator) + dirChild4
// 	if _, err := os.Stat(fullDir4); err == nil {
// 		log.Print(col.RedString("CHEAT"), " - directory /assets/delivery_receipts already exist")
// 	} else if err_create := os.Mkdir(fullDir4, 0755); err_create != nil {
// 		log.Print(col.RedString("CHEAT"), " - error while auto create folder /assets/delivery_receipts")
// 	}

// 	write_image(fullDir1)
// 	write_image(fullDir2)
// 	write_image(fullDir3)
// 	write_image(fullDir4)

// 	log.Print(col.RedString("CHEAT"), " - auto create directory in finished")
// 	return true
// }

// // TODO:
// func write_image(pathUrl string) {
// 	// create blank white image
// 	width := 800
// 	height := 600
// 	img := image.NewRGBA(image.Rect(0, 0, width, height))

// 	// Mengisi citra dengan warna merah
// 	red := color.RGBA{255, 0, 0, 255}
// 	draw.Draw(img, img.Bounds(), &image.Uniform{red}, image.Point{}, draw.Src)

// 	// Membuka file untuk penulisan
// 	file, err := os.Create(pathUrl + string(os.PathSeparator) + "dont-delete.jpg")
// 	if err != nil {
// 		fmt.Println("error while auto create file dont-delete.jpg")
// 	}
// 	defer file.Close()

// 	// Menulis gambar ke dalam format JPG
// 	jpeg.Encode(file, img, nil)
// }
