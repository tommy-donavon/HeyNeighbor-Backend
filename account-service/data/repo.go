package data

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo() *UserRepo {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("DSN"),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
	return &UserRepo{db}
}
