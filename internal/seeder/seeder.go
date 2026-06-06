package seeder

import (
	"fmt"
	"log"
	"math/rand"
	"secureasset-manager/internal/asset"
	"secureasset-manager/internal/cmdb"
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
	// CMDB - Applications
	applications := []cmdb.Application{
		{
			Name:        "SecureAsset Manager",
			Description: "Plateforme ITSM interne de gestion des assets et incidents",
			Version:     "1.0.0",
			Criticality: "High",
			Status:      "Active",
			SiteID:      1,
		},
		{
			Name:        "Portail RH",
			Description: "Application de gestion des demandes RH",
			Version:     "2.3.1",
			Criticality: "Medium",
			Status:      "Active",
			SiteID:      1,
		},
		{
			Name:        "CRM Support",
			Description: "Application utilisée par les équipes support client",
			Version:     "4.8.0",
			Criticality: "Critical",
			Status:      "Active",
			SiteID:      2,
		},
		{
			Name:        "Monitoring Infra",
			Description: "Supervision des serveurs et équipements réseau",
			Version:     "3.1.5",
			Criticality: "Critical",
			Status:      "Active",
			SiteID:      5,
		},
	}

	database.DB.Create(&applications)

	// CMDB - Databases
	databases := []cmdb.Database{
		{
			Name:        "secureasset_prod",
			Engine:      "PostgreSQL",
			Version:     "16",
			Environment: "Production",
			SiteID:      5,
		},
		{
			Name:        "rh_prod",
			Engine:      "MySQL",
			Version:     "8.0",
			Environment: "Production",
			SiteID:      1,
		},
		{
			Name:        "crm_support_prod",
			Engine:      "PostgreSQL",
			Version:     "15",
			Environment: "Production",
			SiteID:      2,
		},
	}

	database.DB.Create(&databases)

	// CMDB - Relations
	relations := []cmdb.ConfigurationRelation{
		{
			SourceType:   "Application",
			SourceID:     1,
			TargetType:   "Asset",
			TargetID:     1,
			RelationType: "hosted_on",
		},
		{
			SourceType:   "Application",
			SourceID:     1,
			TargetType:   "Database",
			TargetID:     1,
			RelationType: "uses_database",
		},
		{
			SourceType:   "Application",
			SourceID:     2,
			TargetType:   "Database",
			TargetID:     2,
			RelationType: "uses_database",
		},
		{
			SourceType:   "Application",
			SourceID:     3,
			TargetType:   "Database",
			TargetID:     3,
			RelationType: "uses_database",
		},
		{
			SourceType:   "Application",
			SourceID:     4,
			TargetType:   "Asset",
			TargetID:     10,
			RelationType: "hosted_on",
		},
	}

	database.DB.Create(&relations)

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

	log.Println("Seeder terminé : 5 sites, 5 services, 72 users, 500 assets, 4 applications, 3 databases, 5 relations CMDB, 10000 incidents")
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
