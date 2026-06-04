package incident

import (
	"time"

	"secureasset-manager/internal/asset"

	"gorm.io/gorm"
)

type Incident struct {
	gorm.Model

	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
	Status      string `json:"status"`

	AssetID uint        `json:"asset_id"`
	Asset   asset.Asset `json:"asset" binding:"-"`

	AssignedTo           string `json:"assigned_to"`
	AssignedTechnicianID uint   `json:"assigned_technician_id"`

	CreatedByID uint `json:"created_by_id"`

	DueAt      *time.Time `json:"due_at"`
	ResolvedAt *time.Time `json:"resolved_at"`
	ClosedAt   *time.Time `json:"closed_at"`
}
