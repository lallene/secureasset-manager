package dashboard

import (
	"net/http"
	"secureasset-manager/internal/asset"
	"time"

	"secureasset-manager/internal/database"
	"secureasset-manager/internal/incident"
	"secureasset-manager/internal/notification"
	"secureasset-manager/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func applyDashboardFilters(query *gorm.DB, siteID string, serviceID string, startDate time.Time) *gorm.DB {
	query = query.Where("created_at >= ?", startDate)

	if siteID != "" {
		query = query.Where("site_id = ?", siteID)
	}

	if serviceID != "" {
		query = query.Where("service_id = ?", serviceID)
	}

	return query
}

func calculateMTTR(priority string, siteID string, serviceID string, startDate time.Time) float64 {
	var mttr float64

	query := applyDashboardFilters(
		database.DB.Model(&incident.Incident{}),
		siteID,
		serviceID,
		startDate,
	)

	query.
		Select("COALESCE(AVG(EXTRACT(EPOCH FROM (resolved_at - created_at)) / 3600), 0)").
		Where("priority = ?", priority).
		Where("resolved_at IS NOT NULL").
		Where("resolved_at > created_at").
		Scan(&mttr)

	return mttr
}

func calculateGlobalMTTR(siteID string, serviceID string, startDate time.Time) float64 {
	var mttr float64

	query := applyDashboardFilters(
		database.DB.Model(&incident.Incident{}),
		siteID,
		serviceID,
		startDate,
	)

	query.
		Select("COALESCE(AVG(EXTRACT(EPOCH FROM (resolved_at - created_at)) / 3600), 0)").
		Where("resolved_at IS NOT NULL").
		Where("resolved_at > created_at").
		Scan(&mttr)

	return mttr
}

func GetStats(c *gin.Context) {
	siteID := c.Query("site_id")
	serviceID := c.Query("service_id")
	period := c.DefaultQuery("period", "30")

	periodDays := 30

	switch period {
	case "7":
		periodDays = 7
	case "30":
		periodDays = 30
	case "90":
		periodDays = 90
	case "180":
		periodDays = 180
	}

	startDate := time.Now().AddDate(0, 0, -periodDays)

	var totalAssets int64
	var activeAssets int64
	var alertAssets int64
	var offlineAssets int64

	var totalIncidents int64
	var openIncidents int64
	var inProgressIncidents int64
	var resolvedIncidents int64
	var closedIncidents int64
	var overdueIncidents int64

	var lowIncidents int64
	var mediumIncidents int64
	var highIncidents int64
	var criticalIncidents int64

	var unreadNotifications int64

	assetQuery := database.DB.Model(&asset.Asset{})

	if siteID != "" {
		assetQuery = assetQuery.Where("site_id = ?", siteID)
	}

	assetQuery.Count(&totalAssets)

	assetQuery.Session(&gorm.Session{}).
		Where("status = ?", "Active").
		Count(&activeAssets)

	assetQuery.Session(&gorm.Session{}).
		Where("status = ?", "Maintenance").
		Count(&alertAssets)

	assetQuery.Session(&gorm.Session{}).
		Where("status = ?", "Retired").
		Count(&offlineAssets)

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Count(&totalIncidents)

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Where("status = ?", "Open").
		Count(&openIncidents)

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Where("status = ?", "In Progress").
		Count(&inProgressIncidents)

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Where("status = ?", "Resolved").
		Count(&resolvedIncidents)

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Where("status = ?", "Closed").
		Count(&closedIncidents)

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Where("priority = ?", "Low").
		Count(&lowIncidents)

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Where("priority = ?", "Medium").
		Count(&mediumIncidents)

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Where("priority = ?", "High").
		Count(&highIncidents)

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Where("priority = ?", "Critical").
		Count(&criticalIncidents)

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Where("status NOT IN ?", []string{"Resolved", "Closed"}).
		Where("due_at IS NOT NULL").
		Where("due_at < ?", time.Now()).
		Count(&overdueIncidents)

	mttrCritical := calculateMTTR("Critical", siteID, serviceID, startDate)
	mttrHigh := calculateMTTR("High", siteID, serviceID, startDate)
	mttrMedium := calculateMTTR("Medium", siteID, serviceID, startDate)
	mttrLow := calculateMTTR("Low", siteID, serviceID, startDate)
	mttrHours := calculateGlobalMTTR(siteID, serviceID, startDate)

	type CategoryMetric struct {
		Type  string `json:"type"`
		Count int64  `json:"count"`
	}

	var categories []CategoryMetric

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Select("type, COUNT(*) as count").
		Group("type").
		Order("count DESC").
		Scan(&categories)

	serviceWorkload := []gin.H{}
	var services []service.Service

	servicesQuery := database.DB.Model(&service.Service{})

	if serviceID != "" {
		servicesQuery = servicesQuery.Where("id = ?", serviceID)
	}

	servicesQuery.Find(&services)

	for _, svc := range services {
		var openCount int64
		var progressCount int64

		applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, "", startDate).
			Where("service_id = ?", svc.ID).
			Where("status = ?", "Open").
			Count(&openCount)

		applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, "", startDate).
			Where("service_id = ?", svc.ID).
			Where("status = ?", "In Progress").
			Count(&progressCount)

		serviceWorkload = append(serviceWorkload, gin.H{
			"service":     svc.Name,
			"open":        openCount,
			"in_progress": progressCount,
		})
	}

	userID, _ := c.Get("user_id")

	if idFloat, ok := userID.(float64); ok {
		database.DB.Model(&notification.Notification{}).
			Where("user_id = ?", uint(idFloat)).
			Where("is_read = ?", false).
			Count(&unreadNotifications)
	}

	weeklyTrendOpened := make([]int64, 7)
	weeklyTrendResolved := make([]int64, 7)

	today := time.Now()

	for i := 0; i < 7; i++ {
		day := today.AddDate(0, 0, -6+i)

		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, day.Location())
		end := start.AddDate(0, 0, 1)

		applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
			Where("created_at >= ? AND created_at < ?", start, end).
			Count(&weeklyTrendOpened[i])

		applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
			Where("resolved_at IS NOT NULL").
			Where("resolved_at >= ? AND resolved_at < ?", start, end).
			Count(&weeklyTrendResolved[i])
	}

	monthlyCreated := make([]int64, 6)
	monthlyResolved := make([]int64, 6)

	for i := 0; i < 6; i++ {
		month := today.AddDate(0, -5+i, 0)

		start := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
		end := start.AddDate(0, 1, 0)

		applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
			Where("created_at >= ? AND created_at < ?", start, end).
			Count(&monthlyCreated[i])

		applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
			Where("resolved_at IS NOT NULL").
			Where("resolved_at >= ? AND resolved_at < ?", start, end).
			Count(&monthlyResolved[i])
	}

	var resolvedCount int64
	var resolvedWithinSLA int64

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Where("resolved_at IS NOT NULL").
		Count(&resolvedCount)

	applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, startDate).
		Where("resolved_at IS NOT NULL").
		Where("due_at IS NOT NULL").
		Where("resolved_at <= due_at").
		Count(&resolvedWithinSLA)

	slaResolutionRate := 0.0

	if resolvedCount > 0 {
		slaResolutionRate = float64(resolvedWithinSLA) / float64(resolvedCount) * 100
	}

	reopeningRate := 0.0

	slaHistory := make([]float64, 12)

	for i := 0; i < 12; i++ {
		start := time.Now().AddDate(0, 0, -7*(12-i))
		end := start.AddDate(0, 0, 7)

		var weekResolved int64
		var weekResolvedWithinSLA int64

		applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, start).
			Where("created_at >= ? AND created_at < ?", start, end).
			Where("resolved_at IS NOT NULL").
			Count(&weekResolved)

		applyDashboardFilters(database.DB.Model(&incident.Incident{}), siteID, serviceID, start).
			Where("created_at >= ? AND created_at < ?", start, end).
			Where("resolved_at IS NOT NULL").
			Where("due_at IS NOT NULL").
			Where("resolved_at <= due_at").
			Count(&weekResolvedWithinSLA)

		if weekResolved > 0 {
			slaHistory[i] = float64(weekResolvedWithinSLA) / float64(weekResolved) * 100
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total_assets":          totalAssets,
		"active_assets":         activeAssets,
		"alert_assets":          alertAssets,
		"offline_assets":        offlineAssets,
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
		"mttr_critical":         mttrCritical,
		"mttr_high":             mttrHigh,
		"mttr_medium":           mttrMedium,
		"mttr_low":              mttrLow,
		"incident_categories":   categories,
		"service_workload":      serviceWorkload,
		"weekly_trend_opened":   weeklyTrendOpened,
		"weekly_trend_resolved": weeklyTrendResolved,
		"monthly_created":       monthlyCreated,
		"monthly_resolved":      monthlyResolved,
		"mttr_hours":            mttrHours,
		"sla_resolution_rate":   slaResolutionRate,
		"reopening_rate":        reopeningRate,
		"sla_history":           slaHistory,
		"filters": gin.H{
			"site_id":    siteID,
			"service_id": serviceID,
			"period":     period,
		},
	})
}

func IsOverdue(item incident.Incident) bool {
	return item.Status != "Closed" &&
		item.Status != "Resolved" &&
		item.DueAt != nil &&
		time.Now().After(*item.DueAt)
}
