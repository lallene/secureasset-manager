package incident

import (
	"net/http"
	"time"

	"secureasset-manager/internal/database"

	"github.com/gin-gonic/gin"
)

func calculateDueDate(priority string) time.Time {
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

func CreateIncident(c *gin.Context) {
	var incident Incident

	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")

	if exists {
		if idFloat, ok := userID.(float64); ok {
			incident.CreatedByID = uint(idFloat)
		}
	}

	incident.Status = "Open"

	dueAt := calculateDueDate(incident.Priority)
	incident.DueAt = &dueAt

	if err := database.DB.Create(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'incident"})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

func GetIncidents(c *gin.Context) {
	var incidents []Incident

	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")

	query := database.DB.Preload("Asset")

	if role == "Viewer" {
		if idFloat, ok := userID.(float64); ok {
			query = query.Where("created_by_id = ?", uint(idFloat))
		}
	}

	if err := query.Find(&incidents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les incidents",
		})
		return
	}

	c.JSON(http.StatusOK, incidents)
}

func GetIncidentByID(c *gin.Context) {
	id := c.Param("id")

	var incident Incident

	if err := database.DB.Preload("Asset").First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Incident introuvable",
		})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func UpdateIncident(c *gin.Context) {
	id := c.Param("id")

	var incident Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Incident introuvable",
		})
		return
	}

	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := database.DB.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de modifier l'incident",
		})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func DeleteIncident(c *gin.Context) {
	id := c.Param("id")

	var incident Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Incident introuvable",
		})
		return
	}

	if err := database.DB.Delete(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de supprimer l'incident",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Incident supprimé avec succès",
	})
}

func TakeIncident(c *gin.Context) {
	id := c.Param("id")

	var incident Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident introuvable"})
		return
	}

	if incident.Status == "Resolved" || incident.Status == "Closed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible de récupérer un ticket résolu ou fermé",
		})
		return
	}

	userID, _ := c.Get("user_id")
	email, _ := c.Get("email")

	if idFloat, ok := userID.(float64); ok {
		incident.AssignedTechnicianID = uint(idFloat)
	}

	if emailString, ok := email.(string); ok {
		incident.AssignedTo = emailString
	}

	incident.Status = "In Progress"

	if err := database.DB.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de prendre le ticket"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func ResolveIncident(c *gin.Context) {
	id := c.Param("id")

	var incident Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident introuvable"})
		return
	}

	if incident.Status != "In Progress" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Seul un ticket en cours peut être résolu",
		})
		return
	}

	now := time.Now()
	incident.Status = "Resolved"
	incident.ResolvedAt = &now

	if err := database.DB.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de résoudre le ticket"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func CloseIncident(c *gin.Context) {
	id := c.Param("id")

	var incident Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident introuvable"})
		return
	}

	if incident.Status != "Resolved" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Seul un ticket résolu peut être clôturé",
		})
		return
	}

	now := time.Now()
	incident.Status = "Closed"
	incident.ClosedAt = &now

	if err := database.DB.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de clôturer le ticket"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

type ReactionInput struct {
	Message string `json:"message" binding:"required"`
}

func ReactToIncident(c *gin.Context) {
	id := c.Param("id")

	var input ReactionInput
	var incident Incident

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident introuvable"})
		return
	}

	if incident.Status == "Closed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible de réagir sur un ticket clôturé",
		})
		return
	}

	incident.Status = "Open"
	incident.ResolvedAt = nil
	incident.ClosedAt = nil

	if err := database.DB.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de réagir au ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Réaction ajoutée, ticket rouvert",
		"reaction": input.Message,
		"incident": incident,
	})
}
