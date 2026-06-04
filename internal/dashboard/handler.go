package dashboard

import (
	"net/http"

	"secureasset-manager/internal/asset"
	"secureasset-manager/internal/database"
	"secureasset-manager/internal/incident"

	"github.com/gin-gonic/gin"
)

func GetStats(c *gin.Context) {
	var totalAssets int64
	var totalIncidents int64
	var openIncidents int64
	var criticalIncidents int64
	var activeAssets int64

	database.DB.Model(&asset.Asset{}).Count(&totalAssets)
	database.DB.Model(&incident.Incident{}).Count(&totalIncidents)
	database.DB.Model(&incident.Incident{}).Where("status = ?", "Open").Count(&openIncidents)
	database.DB.Model(&incident.Incident{}).Where("priority = ?", "Critical").Count(&criticalIncidents)
	database.DB.Model(&asset.Asset{}).Where("status = ?", "Active").Count(&activeAssets)

	c.JSON(http.StatusOK, gin.H{
		"total_assets":       totalAssets,
		"total_incidents":    totalIncidents,
		"open_incidents":     openIncidents,
		"critical_incidents": criticalIncidents,
		"active_assets":      activeAssets,
	})
}
