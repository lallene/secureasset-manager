package incident

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"secureasset-manager/internal/auth"
	"secureasset-manager/internal/database"
	"secureasset-manager/internal/notification"
	ws "secureasset-manager/internal/websocket"

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
	var input Incident

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur non authentifié"})
		return
	}

	var creator auth.User

	if idFloat, ok := userID.(float64); ok {
		if err := database.DB.First(&creator, uint(idFloat)).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Utilisateur introuvable"})
			return
		}
	}

	input.CreatedByID = creator.ID
	input.SiteID = creator.SiteID
	input.Status = "Open"

	dueAt := calculateDueDate(input.Priority)
	input.DueAt = &dueAt

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'incident"})
		return
	}

	var agents []auth.User

	database.DB.
		Where("role = ?", "Agent").
		Where("service_id = ?", input.ServiceID).
		Find(&agents)

	for _, agent := range agents {
		notification.CreateNotification(
			agent.ID,
			"Nouveau ticket",
			input.Title,
			&input.ID,
		)
	}

	ws.Broadcast("Nouveau ticket : " + input.Title)

	c.JSON(http.StatusCreated, input)
}

func GetIncidents(c *gin.Context) {
	var incidents []Incident

	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")

	query := database.DB.
		Preload("Asset").
		Preload("Site").
		Preload("Service")

	if role == "Requester" {
		if idFloat, ok := userID.(float64); ok {
			query = query.Where("created_by_id = ?", uint(idFloat))
		}
	}

	if role == "Agent" {
		var user auth.User

		if idFloat, ok := userID.(float64); ok {
			if err := database.DB.First(&user, uint(idFloat)).Error; err == nil {
				if user.ServiceID != nil {
					query = query.Where("service_id = ?", *user.ServiceID)
				}
			}
		}
	}

	if err := query.Order("created_at DESC").Find(&incidents).Error; err != nil {
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

	if err := database.DB.
		Preload("Asset").
		Preload("Site").
		Preload("Service").
		First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident introuvable"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func UpdateIncident(c *gin.Context) {
	id := c.Param("id")

	var incident Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident introuvable"})
		return
	}

	if incident.Status == "Closed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible de modifier un ticket clôturé",
		})
		return
	}

	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")

	if role == "Requester" {
		if idFloat, ok := userID.(float64); ok {
			if incident.CreatedByID != uint(idFloat) {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Vous ne pouvez modifier que vos propres tickets",
				})
				return
			}
		}
	}

	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if role == "Requester" {
		incident.Status = "Open"
		incident.ResolvedAt = nil
		incident.ClosedAt = nil
	}

	if err := database.DB.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de modifier l'incident"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func DeleteIncident(c *gin.Context) {
	id := c.Param("id")

	var incident Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident introuvable"})
		return
	}

	if incident.Status == "Closed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible de supprimer un ticket clôturé",
		})
		return
	}

	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")

	if role == "Requester" {
		if idFloat, ok := userID.(float64); ok {
			if incident.CreatedByID != uint(idFloat) {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Vous ne pouvez supprimer que vos propres tickets",
				})
				return
			}
		}
	}

	if err := database.DB.Delete(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de supprimer l'incident"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incident supprimé avec succès"})
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

	var agent auth.User

	if idFloat, ok := userID.(float64); ok {
		if err := database.DB.First(&agent, uint(idFloat)).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Agent introuvable"})
			return
		}
	}

	if agent.ServiceID == nil || *agent.ServiceID != incident.ServiceID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Vous ne pouvez prendre que les tickets de votre service",
		})
		return
	}

	incident.AssignedToID = &agent.ID
	incident.Status = "In Progress"

	if err := database.DB.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de prendre le ticket"})
		return
	}

	notification.CreateNotification(
		incident.CreatedByID,
		"Ticket pris en charge",
		"Votre ticket '"+incident.Title+"' est maintenant en cours de traitement.",
		&incident.ID,
	)

	database.DB.Create(&IncidentComment{
		IncidentID: incident.ID,
		UserID:     agent.ID,
		Message:    "Ticket pris en charge",
		Type:       "system",
	})

	ws.Broadcast("Ticket pris en charge : " + incident.Title)

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

	notification.CreateNotification(
		incident.CreatedByID,
		"Ticket résolu",
		"Votre ticket '"+incident.Title+"' a été résolu.",
		&incident.ID,
	)

	userID, _ := c.Get("user_id")

	comment := IncidentComment{
		IncidentID: incident.ID,
		Message:    "Ticket résolu",
		Type:       "system",
	}

	if idFloat, ok := userID.(float64); ok {
		comment.UserID = uint(idFloat)
	}

	database.DB.Create(&comment)

	ws.Broadcast("Ticket résolu : " + incident.Title)

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

	notification.CreateNotification(
		incident.CreatedByID,
		"Ticket clôturé",
		"Votre ticket '"+incident.Title+"' a été clôturé.",
		&incident.ID,
	)

	userID, _ := c.Get("user_id")

	comment := IncidentComment{
		IncidentID: incident.ID,
		Message:    "Ticket clôturé",
		Type:       "system",
	}

	if idFloat, ok := userID.(float64); ok {
		comment.UserID = uint(idFloat)
	}

	database.DB.Create(&comment)

	ws.Broadcast("Ticket clôturé : " + incident.Title)

	c.JSON(http.StatusOK, incident)
}

type CommentInput struct {
	Message string `json:"message" binding:"required"`
}

func AddIncidentComment(c *gin.Context) {
	id := c.Param("id")

	var input CommentInput
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
			"error": "Impossible de commenter un ticket clôturé",
		})
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	comment := IncidentComment{
		IncidentID: incident.ID,
		Message:    input.Message,
		Type:       "comment",
	}

	if idFloat, ok := userID.(float64); ok {
		comment.UserID = uint(idFloat)
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible d'ajouter le commentaire",
		})
		return
	}

	if role == "Requester" {
		incident.Status = "Open"
		incident.ResolvedAt = nil
		incident.ClosedAt = nil
		database.DB.Save(&incident)

		if incident.AssignedToID != nil {
			notification.CreateNotification(
				*incident.AssignedToID,
				"Nouveau commentaire",
				input.Message,
				&incident.ID,
			)
		}
	}

	if role == "Agent" {
		notification.CreateNotification(
			incident.CreatedByID,
			"Nouveau commentaire",
			input.Message,
			&incident.ID,
		)
	}

	ws.Broadcast("Nouveau commentaire sur : " + incident.Title)

	c.JSON(http.StatusCreated, comment)
}

func GetIncidentComments(c *gin.Context) {
	id := c.Param("id")

	var comments []IncidentComment

	if err := database.DB.
		Preload("User").
		Where("incident_id = ?", id).
		Order("created_at ASC").
		Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les commentaires",
		})
		return
	}

	for i := range comments {
		comments[i].User.Password = ""
	}

	c.JSON(http.StatusOK, comments)
}

func UploadIncidentAttachment(c *gin.Context) {
	id := c.Param("id")

	var incident Incident

	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Incident introuvable"})
		return
	}

	if incident.Status == "Closed" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Impossible d'ajouter une pièce jointe sur un ticket clôturé",
		})
		return
	}

	file, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Fichier manquant"})
		return
	}

	uploadDir := "uploads/incidents"

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de créer le dossier upload",
		})
		return
	}

	fileName := fmt.Sprintf("incident_%s_%s", id, file.Filename)
	filePath := filepath.Join(uploadDir, fileName)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible d'enregistrer le fichier",
		})
		return
	}

	userID, _ := c.Get("user_id")

	attachment := IncidentAttachment{
		IncidentID: incident.ID,
		FileName:   file.Filename,
		FilePath:   filePath,
		FileSize:   file.Size,
	}

	if idFloat, ok := userID.(float64); ok {
		attachment.UploadedBy = uint(idFloat)
	}

	if err := database.DB.Create(&attachment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible d'enregistrer la pièce jointe",
		})
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

func GetIncidentAttachments(c *gin.Context) {
	id := c.Param("id")

	var attachments []IncidentAttachment

	if err := database.DB.
		Preload("User").
		Where("incident_id = ?", id).
		Order("created_at DESC").
		Find(&attachments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les pièces jointes",
		})
		return
	}

	for i := range attachments {
		attachments[i].User.Password = ""
	}

	c.JSON(http.StatusOK, attachments)
}

func DownloadIncidentAttachment(c *gin.Context) {
	id := c.Param("id")

	var attachment IncidentAttachment

	if err := database.DB.First(&attachment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Pièce jointe introuvable",
		})
		return
	}

	c.FileAttachment(attachment.FilePath, attachment.FileName)
}

func DeleteIncidentAttachment(c *gin.Context) {
	id := c.Param("id")

	var attachment IncidentAttachment

	if err := database.DB.First(&attachment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Pièce jointe introuvable",
		})
		return
	}

	_ = os.Remove(attachment.FilePath)

	if err := database.DB.Delete(&attachment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de supprimer la pièce jointe",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pièce jointe supprimée",
	})
}
