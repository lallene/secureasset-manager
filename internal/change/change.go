package change

import (
	"secureasset-manager/internal/auth"
	"secureasset-manager/internal/cmdb"
	"time"

	"gorm.io/gorm"
)

type ChangeRequest struct {
	gorm.Model

	Reference string `json:"reference"`

	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`

	Type   string `json:"type"`
	Risk   string `json:"risk"`
	Status string `json:"status"`

	PlannedStart *time.Time `json:"planned_start"`
	PlannedEnd   *time.Time `json:"planned_end"`

	RollbackPlan string `json:"rollback_plan"`

	BusinessServiceID *uint                 `json:"business_service_id"`
	BusinessService   *cmdb.BusinessService `json:"business_service" binding:"-"`

	CreatedByID *uint      `json:"created_by_id"`
	CreatedBy   *auth.User `json:"created_by" binding:"-"`

	ApprovedByID *uint      `json:"approved_by_id"`
	ApprovedBy   *auth.User `json:"approved_by" binding:"-"`

	ApprovedAt    *time.Time `json:"approved_at"`
	RejectedAt    *time.Time `json:"rejected_at"`
	ImplementedAt *time.Time `json:"implemented_at"`
	ClosedAt      *time.Time `json:"closed_at"`
}
