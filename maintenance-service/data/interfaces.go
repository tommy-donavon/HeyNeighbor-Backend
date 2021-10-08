package data

type (
	IMaintenanceRead interface {
		GetPropertyMaintenanceRequest(string) ([]*MaintenanceRequest, error)
		GetMaintenanceRequest(uint) (*MaintenanceRequest, error)
	}

	IMaintenanceCreate interface {
		CreateMaintenanceRequest(*MaintenanceRequest) error
	}

	IMaintenanceDelete interface {
		DeleteMaintenanceRequest(uint) error
	}

	IMaintenancePatch interface {
		UpdateMaintenanceRequest(uint, string, map[string]interface{}) error
	}
)
