package utils

import (
	"io"
	"log"
	"mime/multipart"
	"os"

	"github.com/fatih/color"
)

func Write_image(image multipart.FileHeader, newFileName string, folder string) error {
	// Open the uploaded file
	src, err := image.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	directory := "assets/" + folder + "/"

	// Create a destination file on the server
	dst, err := os.Create(directory + newFileName)
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy the contents of the uploaded file to the destination file
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return err
}

func WriteFile(file multipart.FileHeader, newFileName string, folder string) error {
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	directory := folder

	// Create a destination file on the server
	dst, err := os.Create(directory + newFileName)
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy the contents of the uploaded file to the destination file
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}

	return err
}

func Delete_image(directory string, fileName string) {
	err_delete := os.Remove(directory + fileName)
	if err_delete != nil {
		log.Print(color.RedString("failed to delete image, "), err_delete.Error())
	} else {
		log.Print(color.RedString("success to delete image"))
	}
}
