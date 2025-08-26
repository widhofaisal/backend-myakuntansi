package main

import (
	"backend-file-management/config"
	"backend-file-management/route"
	"backend-file-management/seeder"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func main() {
	config.InitDB()

	e := route.New()

	e.Validator = &CustomValidator{validator: validator.New()}

	// seeder
	seeder.SeedProjects()
	seeder.SeedUser()

	e.Start(":8000")

}
