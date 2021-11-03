package data

import (
	"gorm.io/gorm"
)

type MaintenanceRequest struct {
	gorm.Model
	PropertyCode      string    `json:"property_code" gorm:"not null;index"`
	Tenant            string    `json:"tenant" gorm:"not null"`
	Admin             string    `json:"admin" gorm:"not null"`
	UnitNumber        uint      `json:"unit_number" gorm:"not null"`
	Title             string    `json:"title" gorm:"not null" validate:"required"`
	Description       string    `json:"description" validate:"required"`
	Severity          *severity `json:"severity" validate:"required,gte=0,lte=2"`
	ImageURI          string    `json:"image_uri"`
	InProgress        bool      `json:"in_progress" gorm:"default:0"`
	TenantCheckedDone bool      `json:"tenant_checked_done" gorm:"default:0"`
	AdminCheckDone    bool      `json:"admin_checked_done" gorm:"default:0"`
	IsRejected        bool      `json:"is_rejected" gorm:"default:0"`
	RejectionReason   string    `json:"rejection_reason"`
}

type severity uint

const (
	LOW severity = iota
	MODERATE
	CRITICAL
)
