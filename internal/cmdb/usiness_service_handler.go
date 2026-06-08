package cmdb

import (
	"net/http"

	"secureasset-manager/internal/database"

	"github.com/gin-gonic/gin"
)

func GetBusinessServices(c *gin.Context) {
	var services []BusinessService

	if err := database.DB.
		Preload("Site").
		Preload("Owner").
		Order("created_at DESC").
		Find(&services).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les services métiers"})
		return
	}

	for i := range services {
		if services[i].Owner != nil {
			services[i].Owner.Password = ""
		}
	}

	c.JSON(http.StatusOK, services)
}

func CreateBusinessService(c *gin.Context) {
	var service BusinessService

	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if service.Status == "" {
		service.Status = "Active"
	}

	if err := database.DB.Create(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer le service métier"})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func UpdateBusinessService(c *gin.Context) {
	id := c.Param("id")

	var service BusinessService

	if err := database.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service métier introuvable"})
		return
	}

	if err := c.ShouldBindJSON(&service); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de modifier le service métier"})
		return
	}

	c.JSON(http.StatusOK, service)
}

func DeleteBusinessService(c *gin.Context) {
	id := c.Param("id")

	var service BusinessService

	if err := database.DB.First(&service, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Service métier introuvable"})
		return
	}

	if err := database.DB.Delete(&service).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de supprimer le service métier"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Service métier supprimé avec succès"})
}
