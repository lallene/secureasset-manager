package service

import "gorm.io/gorm"

type Service struct {
	gorm.Model

	Name        string `json:"name" gorm:"unique"`
	Description string `json:"description"`
}
