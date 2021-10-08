package data

import (
	"gorm.io/gorm"
)

type MaintenanceRequest struct {
	gorm.Model
	PropertyCode      string `json:"property_code" gorm:"not null;index"`
	Tenant            string `json:"tenant" gorm:"not null"`
	Admin             string `json:"admin" gorm:"not null"`
	UnitNumber        uint   `json:"unit_number" gorm:"not null"`
	Title             string `json:"title" gorm:"not null" validate:"required"`
	Description       string `json:"description" validate:"required"`
	ImageURI          string `json:"image_uri"`
	TenantCheckedDone bool   `json:"tenant_checked_done" gorm:"default:0"`
	AdminCheckDone    bool   `json:"admin_checked_done" gorm:"default:0"`
	IsRejected        bool   `json:"is_rejected" gorm:"default:0"`
	RejectionReason   string `json:"rejection_reason"`
}
