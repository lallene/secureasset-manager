package cmdb

import (
	"net/http"

	"secureasset-manager/internal/database"

	"github.com/gin-gonic/gin"
)

func GetApplications(c *gin.Context) {
	var applications []Application

	if err := database.DB.
		Preload("Site").
		Preload("Owner").
		Order("created_at DESC").
		Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les applications"})
		return
	}

	for i := range applications {
		if applications[i].Owner != nil {
			applications[i].Owner.Password = ""
		}
	}

	c.JSON(http.StatusOK, applications)
}

func GetApplicationByID(c *gin.Context) {
	id := c.Param("id")

	var application Application

	if err := database.DB.
		Preload("Site").
		Preload("Owner").
		First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application introuvable"})
		return
	}

	if application.Owner != nil {
		application.Owner.Password = ""
	}

	c.JSON(http.StatusOK, application)
}

func CreateApplication(c *gin.Context) {
	var application Application

	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if application.Status == "" {
		application.Status = "Active"
	}

	if err := database.DB.Create(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer l'application"})
		return
	}

	c.JSON(http.StatusCreated, application)
}

func UpdateApplication(c *gin.Context) {
	id := c.Param("id")

	var application Application

	if err := database.DB.First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application introuvable"})
		return
	}

	if err := c.ShouldBindJSON(&application); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de modifier l'application"})
		return
	}

	c.JSON(http.StatusOK, application)
}

func DeleteApplication(c *gin.Context) {
	id := c.Param("id")

	var application Application

	if err := database.DB.First(&application, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Application introuvable"})
		return
	}

	if err := database.DB.Delete(&application).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de supprimer l'application"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Application supprimée avec succès"})
}
