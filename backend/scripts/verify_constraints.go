package main

import (
	"fmt"
	"log"

	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/database"
	"github.com/mo7amedgom3a/arch-visualizer/backend/internal/platform/models"
)

func main() {
	if _, err := database.Connect(); err != nil {
		log.Fatal(err)
	}

	var count int64
	database.DB.Model(&models.ResourceConstraint{}).Count(&count)
	fmt.Printf("Constraints count: %d\n", count)

	var constraints []models.ResourceConstraint
	database.DB.Preload("ResourceType").Find(&constraints)
	for _, c := range constraints {
		fmt.Printf("- %s: %s = %s\n", c.ResourceType.Name, c.ConstraintType, c.ConstraintValue)
	}
}
