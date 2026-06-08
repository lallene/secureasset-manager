package change

import (
	"fmt"
	"net/http"

	"secureasset-manager/internal/database"

	"github.com/gin-gonic/gin"
)

func GetChanges(c *gin.Context) {
	var changes []ChangeRequest

	if err := database.DB.
		Preload("BusinessService").
		Preload("CreatedBy").
		Preload("ApprovedBy").
		Order("created_at DESC").
		Find(&changes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de récupérer les changements"})
		return
	}

	for i := range changes {
		if changes[i].CreatedBy != nil {
			changes[i].CreatedBy.Password = ""
		}
		if changes[i].ApprovedBy != nil {
			changes[i].ApprovedBy.Password = ""
		}
	}

	c.JSON(http.StatusOK, changes)
}

func CreateChange(c *gin.Context) {
	var change ChangeRequest

	if err := c.ShouldBindJSON(&change); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if change.Status == "" {
		change.Status = "Draft"
	}

	if change.Type == "" {
		change.Type = "Standard"
	}

	if change.Risk == "" {
		change.Risk = "Medium"
	}

	if err := database.DB.Create(&change).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer le changement"})
		return
	}

	change.Reference = fmt.Sprintf("CHG-%06d", change.ID)
	database.DB.Save(&change)

	c.JSON(http.StatusCreated, change)
}

func UpdateChange(c *gin.Context) {
	id := c.Param("id")

	var change ChangeRequest

	if err := database.DB.First(&change, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Changement introuvable"})
		return
	}

	if err := c.ShouldBindJSON(&change); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&change).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de modifier le changement"})
		return
	}

	c.JSON(http.StatusOK, change)
}

func DeleteChange(c *gin.Context) {
	id := c.Param("id")

	var change ChangeRequest

	if err := database.DB.First(&change, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Changement introuvable"})
		return
	}

	if err := database.DB.Delete(&change).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de supprimer le changement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Changement supprimé"})
}

func SubmitChange(c *gin.Context) {
	updateChangeStatus(c, "Submitted")
}

func ApproveChange(c *gin.Context) {
	updateChangeStatus(c, "Approved")
}

func RejectChange(c *gin.Context) {
	updateChangeStatus(c, "Rejected")
}

func ImplementChange(c *gin.Context) {
	updateChangeStatus(c, "Implemented")
}

func CloseChange(c *gin.Context) {
	updateChangeStatus(c, "Closed")
}

func updateChangeStatus(c *gin.Context, status string) {
	id := c.Param("id")

	var change ChangeRequest

	if err := database.DB.First(&change, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Changement introuvable"})
		return
	}

	change.Status = status

	if err := database.DB.Save(&change).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de mettre à jour le statut"})
		return
	}

	c.JSON(http.StatusOK, change)
}
