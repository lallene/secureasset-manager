package dashboard

import (
	"net/http"
	"time"

	"secureasset-manager/internal/asset"
	"secureasset-manager/internal/database"
	"secureasset-manager/internal/incident"
	"secureasset-manager/internal/notification"

	"github.com/gin-gonic/gin"
)

func GetStats(c *gin.Context) {
	var totalAssets int64
	var activeAssets int64

	var totalIncidents int64
	var openIncidents int64
	var inProgressIncidents int64
	var resolvedIncidents int64
	var closedIncidents int64
	var criticalIncidents int64
	var overdueIncidents int64

	var unreadNotifications int64
	var lowIncidents int64
	var mediumIncidents int64
	var highIncidents int64

	database.DB.Model(&asset.Asset{}).Count(&totalAssets)
	database.DB.Model(&asset.Asset{}).Where("status = ?", "Active").Count(&activeAssets)

	database.DB.Model(&incident.Incident{}).Count(&totalIncidents)
	database.DB.Model(&incident.Incident{}).Where("status = ?", "Open").Count(&openIncidents)
	database.DB.Model(&incident.Incident{}).Where("status = ?", "In Progress").Count(&inProgressIncidents)
	database.DB.Model(&incident.Incident{}).Where("status = ?", "Resolved").Count(&resolvedIncidents)
	database.DB.Model(&incident.Incident{}).Where("status = ?", "Closed").Count(&closedIncidents)
	database.DB.Model(&incident.Incident{}).Where("priority = ?", "Critical").Count(&criticalIncidents)
	database.DB.Model(&incident.Incident{}).Where("priority = ?", "Low").Count(&lowIncidents)
	database.DB.Model(&incident.Incident{}).Where("priority = ?", "Medium").Count(&mediumIncidents)
	database.DB.Model(&incident.Incident{}).Where("priority = ?", "High").Count(&highIncidents)

	database.DB.Model(&incident.Incident{}).
		Where("status NOT IN ?", []string{"Resolved", "Closed"}).
		Where("due_at IS NOT NULL").
		Where("due_at < ?", time.Now()).
		Count(&overdueIncidents)

	userID, _ := c.Get("user_id")

	if idFloat, ok := userID.(float64); ok {
		database.DB.Model(&notification.Notification{}).
			Where("user_id = ?", uint(idFloat)).
			Where("is_read = ?", false).
			Count(&unreadNotifications)
	}

	c.JSON(http.StatusOK, gin.H{
		"total_assets":          totalAssets,
		"active_assets":         activeAssets,
		"total_incidents":       totalIncidents,
		"open_incidents":        openIncidents,
		"in_progress_incidents": inProgressIncidents,
		"resolved_incidents":    resolvedIncidents,
		"closed_incidents":      closedIncidents,
		"critical_incidents":    criticalIncidents,
		"overdue_incidents":     overdueIncidents,
		"unread_notifications":  unreadNotifications,
		"low_incidents":         lowIncidents,
		"medium_incidents":      mediumIncidents,
		"high_incidents":        highIncidents,
	})
}

func IsOverdue(item incident.Incident) bool {
	return item.Status != "Closed" &&
		item.Status != "Resolved" &&
		item.DueAt != nil &&
		time.Now().After(*item.DueAt)
}
