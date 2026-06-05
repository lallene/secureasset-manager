package notification

import "secureasset-manager/internal/database"

func CreateNotification(
	userID uint,
	title string,
	message string,
	incidentID *uint,
) {
	n := Notification{
		UserID:     userID,
		Title:      title,
		Message:    message,
		IsRead:     false,
		IncidentID: incidentID,
	}

	database.DB.Create(&n)
}
