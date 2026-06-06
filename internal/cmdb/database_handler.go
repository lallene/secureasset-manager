package cmdb

import (
	"net/http"

	"secureasset-manager/internal/database"

	"github.com/gin-gonic/gin"
)

func GetDatabases(c *gin.Context) {
	var databases []Database

	if err := database.DB.
		Preload("Site").
		Order("created_at DESC").
		Find(&databases).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les databases"})
		return
	}

	c.JSON(http.StatusOK, databases)
}

func CreateDatabase(c *gin.Context) {
	var db Database

	if err := c.ShouldBindJSON(&db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&db).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer la database"})
		return
	}

	c.JSON(http.StatusCreated, db)
}

func UpdateDatabase(c *gin.Context) {
	id := c.Param("id")

	var db Database

	if err := database.DB.First(&db, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database introuvable"})
		return
	}

	if err := c.ShouldBindJSON(&db); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&db).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de modifier la database"})
		return
	}

	c.JSON(http.StatusOK, db)
}

func DeleteDatabase(c *gin.Context) {
	id := c.Param("id")

	var db Database

	if err := database.DB.First(&db, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Database introuvable"})
		return
	}

	if err := database.DB.Delete(&db).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de supprimer la database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Database supprimée avec succès"})
}
