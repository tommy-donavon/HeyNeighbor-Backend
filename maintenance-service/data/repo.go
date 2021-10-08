package data

import (
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MaintenanceRepo struct {
	db *gorm.DB
}

var lock = &sync.Mutex{}
var instance *MaintenanceRepo

func NewMaintenanceRepo() *MaintenanceRepo {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			db, err := gorm.Open(postgres.New(postgres.Config{
				DSN:                  os.Getenv("DSN"),
				PreferSimpleProtocol: true,
			}), &gorm.Config{})
			if err != nil {
				panic(err.Error())
			}
			if err := db.AutoMigrate(&MaintenanceRequest{}); err != nil {
				panic(err.Error())
			}
			instance = &MaintenanceRepo{db}
		}
	}
	return instance
}
