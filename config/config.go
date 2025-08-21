package config

import (
	"fmt"
	"os"

	"backend-file-management/model"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	godotenv.Load(".env")
	InitDB()
	InitialMigration()
}

type Config struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_DB       string
}

func InitDB() {
	config := Config{
		DB_Username: os.Getenv("DB_USERNAME"),
		DB_Password: os.Getenv("DB_PASSWORD"),
		DB_Port:     os.Getenv("DB_PORT"),
		DB_Host:     os.Getenv("DB_HOST"),
		DB_DB:       os.Getenv("DB_DB"),
	}

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.DB_Username,
		config.DB_Password,
		config.DB_Host,
		config.DB_Port,
		config.DB_DB,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func InitialMigration() {
	DB.AutoMigrate(
		&model.User{},
		&model.Project{},
		&model.Item{},
		// &model.Shift{},
		// &model.CSS{},
		// &model.CS{},
		// &model.Gerbang{},
		// &model.PicLimNTK{},
		// &model.PicTNIPolri{},
		// &model.BeritaAcaraLimNTK{},
		// &model.BeritaAcaraTNIPolri{},
	)
}
