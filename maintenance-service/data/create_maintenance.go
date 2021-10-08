package data

type IMaintenanceCreate interface {
	CreateMaintenanceRequest(*MaintenanceRequest) error
}

func (mr *MaintenanceRepo) CreateMaintenanceRequest(request *MaintenanceRequest) error {
	return mr.db.Create(request).Error
}
