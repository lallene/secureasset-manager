package incident

import (
	"time"

	"secureasset-manager/internal/asset"
	"secureasset-manager/internal/service"
	"secureasset-manager/internal/site"

	"gorm.io/gorm"
)

type Incident struct {
	gorm.Model

	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`

	Type     string `json:"type" binding:"required"`
	Priority string `json:"priority" binding:"required"`
	Status   string `json:"status"`

	AssetID uint        `json:"asset_id"`
	Asset   asset.Asset `json:"asset" binding:"-"`

	// Site concerné
	SiteID uint      `json:"site_id"`
	Site   site.Site `json:"site" binding:"-"`

	// Service responsable
	ServiceID uint            `json:"service_id"`
	Service   service.Service `json:"service" binding:"-"`

	// Créateur du ticket
	CreatedByID uint `json:"created_by_id"`

	// Agent assigné
	AssignedToID *uint `json:"assigned_to_id"`

	DueAt      *time.Time `json:"due_at"`
	ResolvedAt *time.Time `json:"resolved_at"`
	ClosedAt   *time.Time `json:"closed_at"`
}
