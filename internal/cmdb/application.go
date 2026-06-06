package cmdb

import (
	"secureasset-manager/internal/auth"
	"secureasset-manager/internal/site"

	"gorm.io/gorm"
)

type Application struct {
	gorm.Model

	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Criticality string `json:"criticality" binding:"required"`
	Status      string `json:"status"`

	SiteID uint      `json:"site_id" binding:"required"`
	Site   site.Site `json:"site" binding:"-"`

	OwnerID *uint      `json:"owner_id"`
	Owner   *auth.User `json:"owner" binding:"-"`
}
