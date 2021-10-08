package data

import (
	"encoding/json"
	"fmt"

	my_json "github.com/yhung-mea7/go-rest-kit/data"
)

type IMaintenancePatch interface {
	UpdateMaintenanceRequest(uint, string, map[string]interface{}) error
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
