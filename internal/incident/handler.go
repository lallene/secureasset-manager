package incident

import (
	"net/http"

	"secureasset-manager/internal/database"

	"github.com/gin-gonic/gin"
)

func CreateIncident(c *gin.Context) {
	var incident Incident

	if err := c.ShouldBindJSON(&incident); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'incident"})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

func GetIncidents(c *gin.Context) {
	var incidents []Incident

	if err := database.DB.Preload("Asset").Find(&incidents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les incidents"})
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
