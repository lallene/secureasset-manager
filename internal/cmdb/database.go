package cmdb

import (
	"secureasset-manager/internal/site"

	"gorm.io/gorm"
)

type Database struct {
	gorm.Model

	Name        string `json:"name" binding:"required"`
	Engine      string `json:"engine" binding:"required"`
	Version     string `json:"version"`
	Environment string `json:"environment"`

	SiteID uint      `json:"site_id"`
	Site   site.Site `json:"site" binding:"-"`
}
