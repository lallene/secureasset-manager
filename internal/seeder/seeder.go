package seeder

import (
	"fmt"
	"log"
	"math/rand"
	"secureasset-manager/internal/asset"
	"time"

	"secureasset-manager/internal/auth"
	"secureasset-manager/internal/database"
	"secureasset-manager/internal/incident"
	"secureasset-manager/internal/service"
	"secureasset-manager/internal/site"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed() {
	var count int64
	database.DB.Model(&auth.User{}).Count(&count)

	if count > 0 {
		log.Println("Base déjà alimentée, seeder ignoré")
		return
	}

	log.Println("Seeder nouvelle architecture : génération des données...")

	password, _ := bcrypt.GenerateFromPassword(
		[]byte("password123"),
		bcrypt.DefaultCost,
	)

	// Sites
	sites := []site.Site{
		{Name: "Paris", City: "Paris", Country: "France", Description: "Siège principal"},
		{Name: "Lyon", City: "Lyon", Country: "France", Description: "Site régional"},
		{Name: "Marseille", City: "Marseille", Country: "France", Description: "Site sud"},
		{Name: "Lille", City: "Lille", Country: "France", Description: "Site nord"},
		{Name: "Datacenter", City: "Paris", Country: "France", Description: "Infrastructure centrale"},
	}

	database.DB.Create(&sites)

	// Services
	services := []service.Service{
		{Name: "Réseau", Description: "Gestion LAN/WAN/VPN"},
		{Name: "Systèmes", Description: "Serveurs, AD, GPO, virtualisation"},
		{Name: "Sécurité", Description: "SOC, EDR, vulnérabilités"},
		{Name: "Développement", Description: "Applications internes"},
		{Name: "Support IT", Description: "Support utilisateurs"},
	}

	database.DB.Create(&services)

	var users []auth.User

	// Admins
	for i := 1; i <= 2; i++ {
		siteID := uint(rand.Intn(len(sites)) + 1)

		users = append(users, auth.User{
			Name:     fmt.Sprintf("Admin %02d", i),
			Email:    fmt.Sprintf("admin%d@test.com", i),
			Password: string(password),
			Role:     "Admin",
			SiteID:   siteID,
		})
	}

	// Agents
	for i := 1; i <= 20; i++ {
		siteID := uint(rand.Intn(len(sites)) + 1)
		serviceID := uint(rand.Intn(len(services)) + 1)

		users = append(users, auth.User{
			Name:      fmt.Sprintf("Agent %02d", i),
			Email:     fmt.Sprintf("agent%d@test.com", i),
			Password:  string(password),
			Role:      "Agent",
			SiteID:    siteID,
			ServiceID: &serviceID,
		})
	}

	// Requesters
	for i := 1; i <= 50; i++ {
		siteID := uint(rand.Intn(len(sites)) + 1)

		users = append(users, auth.User{
			Name:     fmt.Sprintf("Requester %02d", i),
			Email:    fmt.Sprintf("requester%d@test.com", i),
			Password: string(password),
			Role:     "Requester",
			SiteID:   siteID,
		})
	}

	database.DB.CreateInBatches(&users, 100)

	// Assets
	var assets []asset.Asset

	assetTypes := []string{"Laptop", "Desktop", "Server", "Printer", "Switch", "Router"}
	assetStatuses := []string{"Active", "Maintenance", "Retired"}

	for i := 1; i <= 500; i++ {
		siteID := uint(rand.Intn(len(sites)) + 1)

		assets = append(assets, asset.Asset{
			Name:       fmt.Sprintf("ASSET-%04d", i),
			Type:       assetTypes[rand.Intn(len(assetTypes))],
			Serial:     fmt.Sprintf("SN-%08d", i),
			IPAddress:  fmt.Sprintf("10.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(254)+1),
			SiteID:     siteID,
			Status:     assetStatuses[rand.Intn(len(assetStatuses))],
			AssignedTo: fmt.Sprintf("Requester %02d", rand.Intn(50)+1),
		})
	}

	database.DB.CreateInBatches(&assets, 500)

	incidentTypes := []string{"Software", "Hardware", "Network", "Security", "Infrastructure"}
	priorities := []string{"Low", "Medium", "High", "Critical"}
	incidentStatuses := []string{"Open", "In Progress", "Resolved", "Closed"}

	var incidents []incident.Incident

	for i := 1; i <= 10000; i++ {
		priority := priorities[rand.Intn(len(priorities))]
		status := incidentStatuses[rand.Intn(len(incidentStatuses))]

		createdAt := time.Now().Add(-time.Duration(rand.Intn(720)+24) * time.Hour)
		updatedAt := createdAt

		dueAt := calculateSeedDueDateFromCreatedAt(priority, createdAt)

		var resolvedAt *time.Time
		var closedAt *time.Time
		var assignedToID *uint

		siteID := uint(rand.Intn(len(sites)) + 1)
		serviceID := uint(rand.Intn(len(services)) + 1)

		requesterID := uint(rand.Intn(50) + 23)

		if status == "In Progress" || status == "Resolved" || status == "Closed" {
			agentID := uint(rand.Intn(20) + 3)
			assignedToID = &agentID
		}

		if status == "Resolved" || status == "Closed" {
			t := createdAt.Add(time.Duration(rand.Intn(72)+1) * time.Hour)
			resolvedAt = &t
			updatedAt = t
		}

		if status == "Closed" && resolvedAt != nil {
			t := resolvedAt.Add(time.Duration(rand.Intn(48)+1) * time.Hour)
			closedAt = &t
			updatedAt = t
		}

		incidents = append(incidents, incident.Incident{
			Model: gorm.Model{
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},

			Title:       fmt.Sprintf("Incident test %05d", i),
			Description: fmt.Sprintf("Description détaillée de l'incident %05d", i),
			Type:        incidentTypes[rand.Intn(len(incidentTypes))],
			Priority:    priority,
			Status:      status,

			AssetID: uint(rand.Intn(500) + 1),

			SiteID:       siteID,
			ServiceID:    serviceID,
			CreatedByID:  requesterID,
			AssignedToID: assignedToID,

			DueAt:      &dueAt,
			ResolvedAt: resolvedAt,
			ClosedAt:   closedAt,
		})
	}

	database.DB.CreateInBatches(&incidents, 500)

	log.Println("Seeder terminé : 5 sites, 5 services, 72 users, 500 assets, 10000 incidents")
}

func calculateSeedDueDateFromCreatedAt(priority string, createdAt time.Time) time.Time {
	switch priority {
	case "Critical":
		return createdAt.Add(4 * time.Hour)
	case "High":
		return createdAt.Add(8 * time.Hour)
	case "Medium":
		return createdAt.Add(24 * time.Hour)
	case "Low":
		return createdAt.Add(72 * time.Hour)
	default:
		return createdAt.Add(24 * time.Hour)
	}
}
