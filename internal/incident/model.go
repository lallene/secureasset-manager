package incident

import (
	"secureasset-manager/internal/asset"

	"gorm.io/gorm"
)

type Incident struct {
	gorm.Model

	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required"`
	Priority    string `json:"priority" binding:"required"`
	Status      string `json:"status" binding:"required"`

	AssetID uint        `json:"asset_id"`
	Asset   asset.Asset `json:"asset" binding:"-"`

	AssignedTo string `json:"assigned_to"`
}
