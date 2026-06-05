package incident

import (
	"secureasset-manager/internal/auth"

	"gorm.io/gorm"
)

type IncidentAttachment struct {
	gorm.Model

	IncidentID uint `json:"incident_id"`
	UploadedBy uint `json:"uploaded_by"`

	FileName string `json:"file_name"`
	FilePath string `json:"file_path"`
	FileSize int64  `json:"file_size"`

	User auth.User `json:"user" gorm:"foreignKey:UploadedBy" binding:"-"`
}
