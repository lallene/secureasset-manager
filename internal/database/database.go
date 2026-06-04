package database

import (
	"fmt"
	"log"
	"time"

	"secureasset-manager/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(cfg config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	var db *gorm.DB
	var err error

	for i := 0; i < 10; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err == nil {
			DB = db
			log.Println("Connexion PostgreSQL réussie")
			return
		}

		log.Println("PostgreSQL indisponible, nouvelle tentative dans 2 secondes...")
		time.Sleep(2 * time.Second)
	}

	log.Fatal("Erreur connexion PostgreSQL :", err)
}
