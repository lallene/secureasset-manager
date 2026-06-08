package seeder

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"secureasset-manager/internal/asset"
	"secureasset-manager/internal/auth"
	"secureasset-manager/internal/cmdb"
	"secureasset-manager/internal/database"
	"secureasset-manager/internal/incident"
	"secureasset-manager/internal/service"
	"secureasset-manager/internal/site"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed() {
	rand.Seed(time.Now().UnixNano())

	var count int64
	database.DB.Model(&auth.User{}).Count(&count)

	if count > 0 {
		log.Println("Base déjà alimentée, seeder ignoré")
		return
	}

	log.Println("Seeder démo ITSM + CMDB : génération des données...")

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Sites
	sites := []site.Site{
		{Name: "Paris HQ", City: "Paris", Country: "France", Description: "Siège principal"},
		{Name: "Lyon Support Center", City: "Lyon", Country: "France", Description: "Centre support IT"},
		{Name: "Marseille Operations", City: "Marseille", Country: "France", Description: "Site opérations sud"},
		{Name: "Lille Contact Center", City: "Lille", Country: "France", Description: "Plateforme utilisateurs"},
		{Name: "Nantes Dev Factory", City: "Nantes", Country: "France", Description: "Centre développement"},
		{Name: "Toulouse Security Lab", City: "Toulouse", Country: "France", Description: "SOC et cybersécurité"},
		{Name: "Datacenter Paris", City: "Paris", Country: "France", Description: "Infrastructure production"},
		{Name: "Datacenter Lyon", City: "Lyon", Country: "France", Description: "Infrastructure PRA"},
	}

	database.DB.Create(&sites)

	// Services
	services := []service.Service{
		{Name: "Réseau", Description: "LAN, WAN, VPN, Wi-Fi"},
		{Name: "Systèmes", Description: "Windows Server, Linux, AD, GPO"},
		{Name: "Cybersécurité", Description: "SOC, EDR, IAM, vulnérabilités"},
		{Name: "Développement", Description: "Applications internes et APIs"},
		{Name: "Support IT", Description: "Support utilisateurs N1/N2"},
		{Name: "Cloud", Description: "Azure, AWS, M365"},
		{Name: "Database", Description: "PostgreSQL, MySQL, sauvegardes"},
		{Name: "Téléphonie", Description: "VoIP, softphones, centres d'appels"},
	}

	database.DB.Create(&services)

	// Users
	var users []auth.User

	for i := 1; i <= 4; i++ {
		users = append(users, auth.User{
			Name:     fmt.Sprintf("Admin %02d", i),
			Email:    fmt.Sprintf("admin%d@test.com", i),
			Password: string(password),
			Role:     "Admin",
			SiteID:   randomID(len(sites)),
		})
	}

	for i := 1; i <= 40; i++ {
		serviceID := randomID(len(services))

		users = append(users, auth.User{
			Name:      fmt.Sprintf("Agent %02d", i),
			Email:     fmt.Sprintf("agent%d@test.com", i),
			Password:  string(password),
			Role:      "Agent",
			SiteID:    randomID(len(sites)),
			ServiceID: &serviceID,
		})
	}

	for i := 1; i <= 60; i++ {
		users = append(users, auth.User{
			Name:     fmt.Sprintf("Requester %02d", i),
			Email:    fmt.Sprintf("requester%d@test.com", i),
			Password: string(password),
			Role:     "Requester",
			SiteID:   randomID(len(sites)),
		})
	}

	database.DB.CreateInBatches(&users, 200)

	// Assets
	var assets []asset.Asset

	assetTypes := []string{"Laptop", "Desktop", "Server", "Printer", "Switch", "Router", "Firewall", "NAS", "Access Point"}
	assetStatuses := []string{"Active", "Active", "Active", "Maintenance", "Retired"}

	for i := 1; i <= 1200; i++ {
		assets = append(assets, asset.Asset{
			Name:       fmt.Sprintf("ASSET-%05d", i),
			Type:       assetTypes[rand.Intn(len(assetTypes))],
			Serial:     fmt.Sprintf("SN-%010d", 1000000000+i),
			IPAddress:  fmt.Sprintf("10.%d.%d.%d", rand.Intn(200)+1, rand.Intn(255), rand.Intn(254)+1),
			SiteID:     randomID(len(sites)),
			Status:     assetStatuses[rand.Intn(len(assetStatuses))],
			AssignedTo: fmt.Sprintf("Requester %02d", rand.Intn(60)+1),
		})
	}

	database.DB.CreateInBatches(&assets, 500)

	// CMDB Applications
	applications := []cmdb.Application{
		{Name: "SecureAsset Manager", Description: "Plateforme ITSM et CMDB interne", Version: "2.0.0", Criticality: "High", Status: "Active", SiteID: 7},
		{Name: "CRM Support", Description: "CRM utilisé par les équipes support client", Version: "4.8.0", Criticality: "Critical", Status: "Active", SiteID: 2},
		{Name: "Portail RH", Description: "Demandes RH et gestion administrative", Version: "2.3.1", Criticality: "Medium", Status: "Active", SiteID: 1},
		{Name: "Monitoring Infra", Description: "Supervision infrastructure et alerting", Version: "3.1.5", Criticality: "Critical", Status: "Active", SiteID: 7},
		{Name: "IAM Portal", Description: "Gestion des identités et accès", Version: "1.9.2", Criticality: "Critical", Status: "Active", SiteID: 6},
		{Name: "Asset Inventory", Description: "Inventaire matériel et logiciel", Version: "5.2.0", Criticality: "High", Status: "Active", SiteID: 7},
		{Name: "Ticketing Legacy", Description: "Ancien outil de ticketing", Version: "8.4.1", Criticality: "Medium", Status: "Maintenance", SiteID: 8},
		{Name: "Finance ERP", Description: "Facturation et comptabilité", Version: "12.1.0", Criticality: "Critical", Status: "Active", SiteID: 1},
		{Name: "Knowledge Base", Description: "Base documentaire IT", Version: "3.0.4", Criticality: "Low", Status: "Active", SiteID: 5},
		{Name: "VPN Gateway Portal", Description: "Portail d'accès VPN utilisateurs", Version: "2.7.8", Criticality: "High", Status: "Active", SiteID: 7},
		{Name: "EDR Console", Description: "Console de supervision EDR", Version: "6.5.2", Criticality: "Critical", Status: "Active", SiteID: 6},
		{Name: "Patch Manager", Description: "Gestion des correctifs systèmes", Version: "4.4.0", Criticality: "High", Status: "Active", SiteID: 7},
		{Name: "M365 Provisioning", Description: "Automatisation des comptes M365", Version: "1.4.3", Criticality: "Medium", Status: "Active", SiteID: 6},
		{Name: "Call Center Tools", Description: "Outils de téléphonie et files d'appels", Version: "9.2.1", Criticality: "High", Status: "Active", SiteID: 4},
		{Name: "Backup Orchestrator", Description: "Orchestration des sauvegardes", Version: "7.0.0", Criticality: "Critical", Status: "Active", SiteID: 8},
		{Name: "DevOps Portal", Description: "CI/CD et déploiements applicatifs", Version: "2.5.6", Criticality: "High", Status: "Active", SiteID: 5},
		{Name: "Data Warehouse", Description: "Reporting BI et extraction données", Version: "10.2.0", Criticality: "High", Status: "Active", SiteID: 7},
		{Name: "Visitor Management", Description: "Gestion des visiteurs", Version: "1.1.0", Criticality: "Low", Status: "Retired", SiteID: 1},
	}

	database.DB.Create(&applications)

	businessServices := []cmdb.BusinessService{
		{
			Name:        "Service ITSM",
			Description: "Gestion incidents, assets et CMDB",
			Criticality: "Critical",
			Status:      "Active",
			SiteID:      1,
		},
		{
			Name:        "Support Client",
			Description: "Support utilisateurs et CRM",
			Criticality: "Critical",
			Status:      "Active",
			SiteID:      2,
		},
		{
			Name:        "Finance & Facturation",
			Description: "ERP financier",
			Criticality: "Critical",
			Status:      "Active",
			SiteID:      1,
		},
		{
			Name:        "Cyber Security Operations",
			Description: "SOC, IAM, EDR",
			Criticality: "Critical",
			Status:      "Active",
			SiteID:      6,
		},
		{
			Name:        "Infrastructure & Monitoring",
			Description: "Datacenters et supervision",
			Criticality: "High",
			Status:      "Active",
			SiteID:      7,
		},
		{
			Name:        "Digital Workplace",
			Description: "M365, VPN, postes utilisateurs",
			Criticality: "High",
			Status:      "Active",
			SiteID:      1,
		},
	}

	database.DB.Create(&businessServices)

	// CMDB Databases
	databases := []cmdb.Database{
		{Name: "secureasset_prod", Engine: "PostgreSQL", Version: "16", Environment: "Production", SiteID: 7},
		{Name: "secureasset_pra", Engine: "PostgreSQL", Version: "16", Environment: "Preproduction", SiteID: 8},
		{Name: "crm_support_prod", Engine: "PostgreSQL", Version: "15", Environment: "Production", SiteID: 2},
		{Name: "rh_prod", Engine: "MySQL", Version: "8.0", Environment: "Production", SiteID: 1},
		{Name: "iam_prod", Engine: "PostgreSQL", Version: "15", Environment: "Production", SiteID: 6},
		{Name: "finance_erp_prod", Engine: "Oracle", Version: "19c", Environment: "Production", SiteID: 1},
		{Name: "monitoring_tsdb", Engine: "PostgreSQL", Version: "14", Environment: "Production", SiteID: 7},
		{Name: "knowledge_base_prod", Engine: "MySQL", Version: "8.0", Environment: "Production", SiteID: 5},
		{Name: "edr_events", Engine: "MongoDB", Version: "6", Environment: "Production", SiteID: 6},
		{Name: "patch_manager_prod", Engine: "SQL Server", Version: "2022", Environment: "Production", SiteID: 7},
		{Name: "callcenter_prod", Engine: "PostgreSQL", Version: "15", Environment: "Production", SiteID: 4},
		{Name: "backup_catalog", Engine: "PostgreSQL", Version: "15", Environment: "Production", SiteID: 8},
		{Name: "devops_registry", Engine: "Redis", Version: "7", Environment: "Production", SiteID: 5},
		{Name: "bi_warehouse", Engine: "PostgreSQL", Version: "16", Environment: "Production", SiteID: 7},
	}

	database.DB.Create(&databases)

	// CMDB Relations
	var relations []cmdb.ConfigurationRelation

	for i := 1; i <= len(applications); i++ {
		relations = append(relations, cmdb.ConfigurationRelation{
			SourceType:   "Application",
			SourceID:     uint(i),
			TargetType:   "Asset",
			TargetID:     uint(rand.Intn(1200) + 1),
			RelationType: "hosted_on",
		})
	}

	for i := 1; i <= len(applications); i++ {
		dbID := uint((i-1)%len(databases) + 1)

		relations = append(relations, cmdb.ConfigurationRelation{
			SourceType:   "Application",
			SourceID:     uint(i),
			TargetType:   "Database",
			TargetID:     dbID,
			RelationType: "uses_database",
		})
	}

	for i := 1; i <= 9; i++ {
		relations = append(relations, cmdb.ConfigurationRelation{
			SourceType:   "Application",
			SourceID:     uint(rand.Intn(len(applications)) + 1),
			TargetType:   "Application",
			TargetID:     uint(rand.Intn(len(applications)) + 1),
			RelationType: "depends_on",
		})
	}

	relations = append(relations,

		cmdb.ConfigurationRelation{
			SourceType:   "BusinessService",
			SourceID:     1,
			TargetType:   "Application",
			TargetID:     1,
			RelationType: "depends_on",
		},

		cmdb.ConfigurationRelation{
			SourceType:   "BusinessService",
			SourceID:     2,
			TargetType:   "Application",
			TargetID:     2,
			RelationType: "depends_on",
		},

		cmdb.ConfigurationRelation{
			SourceType:   "BusinessService",
			SourceID:     3,
			TargetType:   "Application",
			TargetID:     8,
			RelationType: "depends_on",
		},

		cmdb.ConfigurationRelation{
			SourceType:   "BusinessService",
			SourceID:     4,
			TargetType:   "Application",
			TargetID:     5,
			RelationType: "depends_on",
		},

		cmdb.ConfigurationRelation{
			SourceType:   "BusinessService",
			SourceID:     4,
			TargetType:   "Application",
			TargetID:     11,
			RelationType: "depends_on",
		},

		cmdb.ConfigurationRelation{
			SourceType:   "BusinessService",
			SourceID:     5,
			TargetType:   "Application",
			TargetID:     4,
			RelationType: "depends_on",
		},
	)
	database.DB.Create(&relations)

	// Incidents
	incidentTypes := []string{"Software", "Hardware", "Network", "Security", "Infrastructure", "Access", "Database", "Cloud"}
	priorities := []string{"Low", "Medium", "High", "Critical"}
	incidentStatuses := []string{"Open", "In Progress", "Resolved", "Closed"}

	var incidents []incident.Incident

	for i := 1; i <= 20000; i++ {
		priority := priorities[rand.Intn(len(priorities))]
		status := incidentStatuses[rand.Intn(len(incidentStatuses))]

		createdAt := time.Now().Add(-time.Duration(rand.Intn(1440)+12) * time.Hour)
		updatedAt := createdAt
		dueAt := calculateSeedDueDateFromCreatedAt(priority, createdAt)

		var resolvedAt *time.Time
		var closedAt *time.Time
		var assignedToID *uint

		requesterID := uint(rand.Intn(60) + 45)

		if status == "In Progress" || status == "Resolved" || status == "Closed" {
			agentID := uint(rand.Intn(40) + 5)
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

			Title:       randomIncidentTitle(i),
			Description: randomIncidentDescription(),
			Type:        incidentTypes[rand.Intn(len(incidentTypes))],
			Priority:    priority,
			Status:      status,

			AssetID: uint(rand.Intn(1200) + 1),

			SiteID:       randomID(len(sites)),
			ServiceID:    randomID(len(services)),
			CreatedByID:  requesterID,
			AssignedToID: assignedToID,

			DueAt:      &dueAt,
			ResolvedAt: resolvedAt,
			ClosedAt:   closedAt,
		})
	}

	database.DB.CreateInBatches(&incidents, 500)

	log.Println("Seeder terminé : 8 sites, 8 services, 104 users, 1200 assets, 18 applications, 14 databases, 51 relations CMDB, 20000 incidents")
}

func randomID(max int) uint {
	return uint(rand.Intn(max) + 1)
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

func randomIncidentTitle(i int) string {
	titles := []string{
		"Perte de connectivité VPN",
		"Poste utilisateur lent",
		"Erreur authentification M365",
		"Serveur applicatif indisponible",
		"Alerte EDR sur endpoint",
		"Imprimante réseau inaccessible",
		"Erreur sauvegarde nocturne",
		"Latence applicative élevée",
		"Compte utilisateur verrouillé",
		"Base de données saturée",
		"Certificat SSL expiré",
		"Incident Wi-Fi étage 3",
		"Patch critique en échec",
		"Suspicion phishing utilisateur",
		"Déploiement applicatif échoué",
	}

	return fmt.Sprintf("%s #%05d", titles[rand.Intn(len(titles))], i)
}

func randomIncidentDescription() string {
	descriptions := []string{
		"Un utilisateur signale une anomalie nécessitant une analyse technique.",
		"Le monitoring a généré une alerte nécessitant une vérification.",
		"L'incident impacte un service métier et doit être priorisé.",
		"Une action corrective est attendue par l'équipe assignée.",
		"Les premiers éléments indiquent une indisponibilité partielle.",
	}

	return descriptions[rand.Intn(len(descriptions))]
}
