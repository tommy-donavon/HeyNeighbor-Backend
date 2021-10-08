package data

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MaintenanceRepo struct {
	db *gorm.DB
}

func NewMaintenanceRepo() *MaintenanceRepo {
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
	return &MaintenanceRepo{db}
}
