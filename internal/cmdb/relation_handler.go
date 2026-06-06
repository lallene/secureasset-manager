package cmdb

import (
	"net/http"

	"secureasset-manager/internal/database"

	"github.com/gin-gonic/gin"
)

func GetRelations(c *gin.Context) {
	var relations []ConfigurationRelation

	if err := database.DB.
		Order("created_at DESC").
		Find(&relations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les relations"})
		return
	}

	c.JSON(http.StatusOK, relations)
}

func CreateRelation(c *gin.Context) {
	var relation ConfigurationRelation

	if err := c.ShouldBindJSON(&relation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&relation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer la relation"})
		return
	}

	c.JSON(http.StatusCreated, relation)
}

func UpdateRelation(c *gin.Context) {
	id := c.Param("id")

	var relation ConfigurationRelation

	if err := database.DB.First(&relation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Relation introuvable"})
		return
	}

	if err := c.ShouldBindJSON(&relation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&relation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de modifier la relation"})
		return
	}

	c.JSON(http.StatusOK, relation)
}

func DeleteRelation(c *gin.Context) {
	id := c.Param("id")

	var relation ConfigurationRelation

	if err := database.DB.First(&relation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Relation introuvable"})
		return
	}

	if err := database.DB.Delete(&relation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de supprimer la relation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Relation supprimée avec succès"})
}
