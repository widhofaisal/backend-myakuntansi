package seeder

import (
	"backend-file-management/config"
	"backend-file-management/model"
	"log"
)

func SeedUser() {
	var count int64
	config.DB.Model(&model.User{}).Count(&count)
	if count == 0 {
		users := []model.User{
			{
				Fullname: "Adam Hikmawan",
				Username: "adamh",
				Password: "$2y$14$CktAhvc4iEl6XG0ej7qZY.nrblR.sCxkEAEu655ehBvCqI/BK/HsS", // qwerty123
				Role:     "admin",
			},
			{
				Fullname: "widho Faisal Hakim",
				Username: "widhofh",
				Password: "$2y$14$CktAhvc4iEl6XG0ej7qZY.nrblR.sCxkEAEu655ehBvCqI/BK/HsS", // qwerty123
				Role:     "user",
			},
		}

		for _, user := range users {
			if err := config.DB.Create(&user).Error; err != nil {
				log.Printf("Seeder error (user): %v", err)
			}
		}
		log.Println("✅ SeedUser completed.")
	}
	log.Println("✅ SeedUser already exists.")
}
