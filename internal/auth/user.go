package auth

import (
	"secureasset-manager/internal/service"
	"secureasset-manager/internal/site"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" gorm:"unique" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required"` // Admin | Agent | Requester

	SiteID    uint  `json:"site_id"`
	ServiceID *uint `json:"service_id"`

	Site    site.Site       `json:"site" binding:"-"`
	Service service.Service `json:"service" binding:"-"`
}
