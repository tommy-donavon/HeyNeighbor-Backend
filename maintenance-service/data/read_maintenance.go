package data

type IMaintenanceRead interface {
	GetPropertyMaintenanceRequest(string) ([]*MaintenanceRequest, error)
	GetMaintenanceRequest(uint) (*MaintenanceRequest, error)
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
