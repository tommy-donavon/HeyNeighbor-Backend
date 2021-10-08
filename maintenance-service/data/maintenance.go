package data

import (
	"encoding/json"
	"fmt"

	my_json "github.com/yhung-mea7/go-rest-kit/data"
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

func (mr *MaintenanceRepo) GetPropertyMaintenanceRequest(propertyCode string) ([]*MaintenanceRequest, error) {
	request := []*MaintenanceRequest{}
	if err := mr.db.Where("property_code = ?", propertyCode).Find(&request).Error; err != nil {
		return nil, err
	}
	return request, nil
}

func (mr *MaintenanceRepo) GetMaintenanceRequest(id uint) (*MaintenanceRequest, error) {
	request := MaintenanceRequest{}
	if err := mr.db.Where("id = ?", id).First(&request).Error; err != nil {
		return nil, err
	}
	return &request, nil
}

func (mr *MaintenanceRepo) CreateMaintenanceRequest(request *MaintenanceRequest) error {
	return mr.db.Create(request).Error
}

func (mr *MaintenanceRepo) DeleteMaintenanceRequest(id uint) error {
	request, err := mr.GetMaintenanceRequest(id)
	if err != nil {
		return err
	}
	return mr.db.Delete(&request).Error
}

func (mr *MaintenanceRepo) UpdateMaintenanceRequest(id uint, username string, updateInfo map[string]interface{}) error {
	request, err := mr.GetMaintenanceRequest(id)
	if err != nil {
		return err
	}

	if request.Admin != username && request.Tenant != username {
		return fmt.Errorf("user is not allowed to alter this request")
	}
	rBytes, err := json.Marshal(request)
	if err != nil {
		return err
	}
	requestMap := map[string]interface{}{}
	if err := json.Unmarshal(rBytes, &requestMap); err != nil {
		return err
	}
	for key, value := range updateInfo {
		if _, ok := requestMap[key]; ok {
			switch key {
			case "title":
				v, ok := value.(string)
				if !ok {
					return fmt.Errorf("%s can not be asserted to string", value)
				}
				request.Title = v
			case "description":
				v, ok := value.(string)
				if !ok {
					return fmt.Errorf("%s can not be asserted to string", value)
				}
				request.Description = v
			case "image_uri":
				v, ok := value.(string)
				if !ok {
					return fmt.Errorf("%s can not be asserted to string", value)
				}
				request.ImageURI = v
			case "tenant_checked_done":
				v, ok := value.(bool)
				if !ok {
					return fmt.Errorf("%s can not be asserted to a boolean", value)
				}
				request.TenantCheckedDone = v
			case "admin_checked_done":
				v, ok := value.(bool)
				if !ok {
					return fmt.Errorf("%s can not be asserted to a boolean", value)
				}
				request.AdminCheckDone = v
			case "is_rejected":
				v, ok := value.(bool)
				if !ok {
					return fmt.Errorf("%s can not be asserted to a boolean", value)
				}
				request.IsRejected = v
			case "rejection_reason":
				v, ok := value.(string)
				if !ok {
					return fmt.Errorf("%s can not be asserted to string", value)
				}
				request.RejectionReason = v
			default:
				return fmt.Errorf("%s is not an updatable field", key)
			}
		}
	}
	if err := my_json.NewValidator().Validate(request); err != nil {
		return err
	}
	return mr.db.Save(&request).Error

}
