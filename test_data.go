package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/database"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/models"
	"github.com/malyyboh/slowo-wiary-warszawa-bot/internal/repository"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if err := database.InitDB("./data/bot.db"); err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer database.CloseDB()

	repo := repository.NewEventRepository()

	events := []models.Event{
		{
			Title:       "Різдвяне богослужіння",
			Description: "Святкове різдвяне богослужіння з особливою програмою та прославленням",
			Date:        time.Date(2025, 12, 25, 16, 0, 0, 0, time.UTC),
			Location:    stringPtr("Wyborna 20, Warszawa"),
			Category:    stringPtr("Богослужіння"),
			IsPublished: true,
			CreatedAt:   time.Now(),
			CreatedBy:   123456789,
		},
		{
			Title:           "Молодіжний семінар",
			Description:     "Семінар для молоді на тему \"Віра в сучасному світі\"",
			Date:            time.Date(2025, 11, 30, 18, 0, 0, 0, time.UTC),
			Location:        stringPtr("Trakt Lubelski 410, Warszawa"),
			Category:        stringPtr("Семінар"),
			RegistrationURL: stringPtr("https://forms.google.com/example"),
			IsPublished:     true,
			CreatedAt:       time.Now(),
			CreatedBy:       123456789,
		},
		{
			Title:       "Молитовна зустріч",
			Description: "Спеціальна молитовна зустріч за мир в Україні",
			Date:        time.Date(2025, 11, 20, 19, 0, 0, 0, time.UTC),
			IsPublished: true,
			CreatedAt:   time.Now(),
			CreatedBy:   123456789,
		},
	}

	for _, event := range events {
		if err := repo.Create(&event); err != nil {
			log.Printf("Error creating event: %v", err)
		} else {
			log.Printf("✅ Created event: %s (ID: %d)", event.Title, event.ID)
		}
	}

	log.Println("✅ Test data created successfully!")
}

func stringPtr(s string) *string {
	return &s
}
