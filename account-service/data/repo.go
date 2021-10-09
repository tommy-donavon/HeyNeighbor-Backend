package data

import (
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

var lock = &sync.Mutex{}
var repoInstance *UserRepo

func NewUserRepo() *UserRepo {
	if repoInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if repoInstance == nil {
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
			repoInstance = &UserRepo{db}
		}
	}
	return repoInstance
}
