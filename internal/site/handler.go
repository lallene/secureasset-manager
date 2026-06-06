package site

import (
	"net/http"

	"secureasset-manager/internal/database"

	"github.com/gin-gonic/gin"
)

func GetSites(c *gin.Context) {
	var sites []Site

	if err := database.DB.Order("name ASC").Find(&sites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les sites"})
		return
	}

	c.JSON(http.StatusOK, sites)
}

func CreateSite(c *gin.Context) {
	var site Site

	if err := c.ShouldBindJSON(&site); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer le site"})
		return
	}

	c.JSON(http.StatusCreated, site)
}

func UpdateSite(c *gin.Context) {
	id := c.Param("id")

	var site Site

	if err := database.DB.First(&site, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site introuvable"})
		return
	}

	if err := c.ShouldBindJSON(&site); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de modifier le site"})
		return
	}

	c.JSON(http.StatusOK, site)
}

func DeleteSite(c *gin.Context) {
	id := c.Param("id")

	var site Site

	if err := database.DB.First(&site, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Site introuvable"})
		return
	}

	if err := database.DB.Delete(&site).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de supprimer le site"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Site supprimé avec succès"})
}
