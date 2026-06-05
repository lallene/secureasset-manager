package notification

import (
	"secureasset-manager/internal/auth"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model

	UserID uint `json:"user_id"`

	Title   string `json:"title"`
	Message string `json:"message"`
	IsRead  bool   `json:"is_read"`

	IncidentID *uint `json:"incident_id"`

	User auth.User `json:"user" binding:"-"`
}
