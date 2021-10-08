package data

type IMaintenanceDelete interface {
	DeleteMaintenanceRequest(uint) error
}

func (mr *MaintenanceRepo) DeleteMaintenanceRequest(id uint) error {
	request, err := mr.GetMaintenanceRequest(id)
	if err != nil {
		return err
	}
	return mr.db.Delete(&request).Error
}
