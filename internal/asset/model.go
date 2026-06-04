package asset

import "gorm.io/gorm"

type Asset struct {
	gorm.Model

	Name string `json:"name" binding:"required"`
	Type string `json:"type" binding:"required"`

	Serial string `json:"serial"`

	IPAddress string `json:"ip_address" binding:"required,ip"`

	Site string `json:"site" binding:"required"`

	Status string `json:"status" binding:"required"`

	AssignedTo string `json:"assigned_to"`
}
