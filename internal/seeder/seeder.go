package seeder

import (
	"fmt"
	"log"
	"math/rand"
	"secureasset-manager/internal/database"
	"time"

	"secureasset-manager/internal/asset"
	"secureasset-manager/internal/auth"
	"secureasset-manager/internal/incident"

	"golang.org/x/crypto/bcrypt"
)

func Seed() {
	var count int64
	database.DB.Model(&auth.User{}).Count(&count)

	if count > 0 {
		log.Println("Base déjà alimentée, seeder ignoré")
		return
	}

	log.Println("Seeder : génération des données de test...")

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	roles := []string{"Admin", "Technician", "Viewer"}

	users := []auth.User{
		{Name: "Cedric Admin", Email: "admin@test.com", Password: string(password), Role: "Admin"},
		{Name: "Technicien Réseau", Email: "tech1@test.com", Password: string(password), Role: "Technician"},
		{Name: "Technicien Support", Email: "tech2@test.com", Password: string(password), Role: "Technician"},
		{Name: "Utilisateur Viewer", Email: "viewer@test.com", Password: string(password), Role: "Viewer"},
	}

	for i := 5; i <= 20; i++ {
		role := roles[rand.Intn(len(roles))]

		users = append(users, auth.User{
			Name:     fmt.Sprintf("User Test %02d", i),
			Email:    fmt.Sprintf("user%02d@test.com", i),
			Password: string(password),
			Role:     role,
		})
	}

	database.DB.Create(&users)

	assetTypes := []string{"Laptop", "Desktop", "Server", "Printer", "Switch", "Router"}
	sites := []string{"Paris", "Lyon", "Marseille", "Lille", "Datacenter"}
	statuses := []string{"Active", "Maintenance", "Retired"}

	var assets []asset.Asset

	for i := 1; i <= 50; i++ {
		assets = append(assets, asset.Asset{
			Name:       fmt.Sprintf("ASSET-%03d", i),
			Type:       assetTypes[rand.Intn(len(assetTypes))],
			Serial:     fmt.Sprintf("SN-%06d", i),
			IPAddress:  fmt.Sprintf("192.168.%d.%d", rand.Intn(10)+1, rand.Intn(200)+10),
			Site:       sites[rand.Intn(len(sites))],
			Status:     statuses[rand.Intn(len(statuses))],
			AssignedTo: fmt.Sprintf("User Test %02d", rand.Intn(20)+1),
		})
	}

	database.DB.Create(&assets)

	incidentTypes := []string{"Software", "Hardware", "Network", "Security", "Infrastructure"}
	priorities := []string{"Low", "Medium", "High", "Critical"}
	incidentStatuses := []string{"Open", "In Progress", "Resolved", "Closed"}

	var incidents []incident.Incident

	for i := 1; i <= 100; i++ {
		priority := priorities[rand.Intn(len(priorities))]
		status := incidentStatuses[rand.Intn(len(incidentStatuses))]

		dueAt := calculateSeedDueDate(priority)

		var resolvedAt *time.Time
		var closedAt *time.Time

		if status == "Resolved" || status == "Closed" {
			t := time.Now().Add(-time.Duration(rand.Intn(48)) * time.Hour)
			resolvedAt = &t
		}

		if status == "Closed" {
			t := time.Now().Add(-time.Duration(rand.Intn(24)) * time.Hour)
			closedAt = &t
		}

		incidents = append(incidents, incident.Incident{
			Title:       fmt.Sprintf("Incident test %03d", i),
			Description: fmt.Sprintf("Description détaillée de l'incident %03d", i),
			Type:        incidentTypes[rand.Intn(len(incidentTypes))],
			Priority:    priority,
			Status:      status,

			AssetID: uint(rand.Intn(50) + 1),

			CreatedByID:          uint(rand.Intn(20) + 1),
			AssignedTechnicianID: uint(rand.Intn(2) + 2),
			AssignedTo:           fmt.Sprintf("tech%d@test.com", rand.Intn(2)+1),

			DueAt:      &dueAt,
			ResolvedAt: resolvedAt,
			ClosedAt:   closedAt,
		})
	}

	database.DB.Create(&incidents)

	log.Println("Seeder terminé : 20 users, 50 assets, 100 incidents")
}

func calculateSeedDueDate(priority string) time.Time {
	now := time.Now()

	switch priority {
	case "Critical":
		return now.Add(4 * time.Hour)
	case "High":
		return now.Add(8 * time.Hour)
	case "Medium":
		return now.Add(24 * time.Hour)
	case "Low":
		return now.Add(72 * time.Hour)
	default:
		return now.Add(24 * time.Hour)
	}
}
