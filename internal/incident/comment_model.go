package incident

import (
	"secureasset-manager/internal/auth"

	"gorm.io/gorm"
)

type IncidentComment struct {
	gorm.Model

	IncidentID uint `json:"incident_id"`
	UserID     uint `json:"user_id"`

	Message string `json:"message" binding:"required"`
	Type    string `json:"type"`

	User auth.User `json:"user" binding:"-"`
}
