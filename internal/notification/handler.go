package notification

import (
	"net/http"

	"secureasset-manager/internal/database"

	"github.com/gin-gonic/gin"
)

func GetNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var notifications []Notification

	query := database.DB.Order("created_at DESC")

	if idFloat, ok := userID.(float64); ok {
		query = query.Where("user_id = ?", uint(idFloat))
	}

	if err := query.Find(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de récupérer les notifications",
		})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

func MarkAsRead(c *gin.Context) {
	id := c.Param("id")

	var notification Notification

	if err := database.DB.First(&notification, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Notification introuvable",
		})
		return
	}

	notification.IsRead = true

	if err := database.DB.Save(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Impossible de marquer comme lu",
		})
		return
	}

	c.JSON(http.StatusOK, notification)
}
