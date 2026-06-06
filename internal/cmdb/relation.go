package cmdb

import "gorm.io/gorm"

type ConfigurationRelation struct {
	gorm.Model

	SourceType string `json:"source_type" binding:"required"`
	SourceID   uint   `json:"source_id" binding:"required"`

	TargetType string `json:"target_type" binding:"required"`
	TargetID   uint   `json:"target_id" binding:"required"`

	RelationType string `json:"relation_type" binding:"required"`
}
