package asset

import (
	"secureasset-manager/internal/site"
	"time"

	"gorm.io/gorm"
)

type Asset struct {
	gorm.Model

	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`

	Serial string `json:"serial"`

	IPAddress string `json:"ip_address" binding:"required,ip"`

	SiteID uint      `json:"site_id"`
	Site   site.Site `json:"site"`

	Status string `json:"status" binding:"required"`

	AssignedTo string `json:"assigned_to"`

	AssignedTechnicianID uint       `json:"assigned_technician_id"`
	DueAt                *time.Time `json:"due_at"`
	ResolvedAt           *time.Time `json:"resolved_at"`
	ClosedAt             *time.Time `json:"closed_at"`
}
