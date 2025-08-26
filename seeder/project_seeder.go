package seeder

import (
	"backend-file-management/config"
	"backend-file-management/model"
	"log"
)

func SeedProjects() {
	var count int64
	config.DB.Model(&model.Project{}).Count(&count)
	if count == 0 {
		projects := []model.Project{
			{
				Name:        "E-Government Archive System",
				Description: "A document system for internal use",
			},
			{
				Name:        "Digital Letter Management",
				Description: "Incoming/outgoing letter tracker",
			},
			{
				Name:        "Internal Cloud Storage",
				Description: "File storage and sharing for teams",
			},
		}

		for _, project := range projects {
			if err := config.DB.Create(&project).Error; err != nil {
				log.Printf("Seeder error (project): %v", err)
			}
		}
		log.Println("✅ SeedProjects completed.")
	}
	log.Println("✅ SeedProjects already exists.")
}
