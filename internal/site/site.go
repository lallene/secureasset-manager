package site

import "gorm.io/gorm"

type Site struct {
	gorm.Model

	Name        string `json:"name" gorm:"unique"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Description string `json:"description"`
}
