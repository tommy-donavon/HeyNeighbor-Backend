package data

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		Username    string `json:"username" validate:"required" gorm:"unique"`
		Password    string `json:"password" validate:"required"`
		FirstName   string `json:"first_name" validate:"required"`
		LastName    string `json:"last_name" validate:"required"`
		UnitNumber  uint   `json:"unit_number" validate:"required"`
		Email       string `json:"email" validate:"email" gorm:"unique"`
		PhoneNumber string `json:"phone_number" validate:"phone" gorm:"unique"`
		ProfileURI  string `json:"profile_uri"`
		AccountType uint   `json:"account_type" validate:"gte=0,lte=1"`
	}

	UserRepo struct {
		db *gorm.DB
	}
)

func NewUserRepo() *UserRepo {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("DSN"),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&User{})
	return &UserRepo{db}
}

func (ur *UserRepo) CreateUser(user *User) error {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash
	return ur.db.Create(user).Error
}

func (ur *UserRepo) GetUser(id uint) (*User, error) {
	user := User{}
	err := ur.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (ur *UserRepo) UpdateUser(id uint, updateInfo map[string]string) error {
	user, err := ur.GetUser(id)
	if err != nil {
		return err
	}
	userInfoBytes, err := json.Marshal(user)
	if err != nil {
		return err
	}
	userMap := map[string]interface{}{}
	err = json.Unmarshal(userInfoBytes, &userMap)
	if err != nil {
		return err
	}
	for key, value := range updateInfo {
		if _, ok := userMap[key]; ok {
			switch key {
			case "username":
				user.Username = value
			case "password":
				hashP, err := hashPassword(value)
				if err != nil {
					return err
				}
				user.Password = hashP
			case "first_name":
				user.FirstName = value
			case "last_name":
				user.LastName = value
			case "unit_number":
				u, err := strconv.ParseUint(value, 10, 32)
				if err != nil {
					return err
				}
				user.UnitNumber = uint(u)
			case "email":
				user.Email = value
			case "phone_number":
				user.PhoneNumber = value
			case "profile_uri":
				user.ProfileURI = value
			}
		} else {
			return fmt.Errorf("%s is not a valid field", key)
		}

	}

	if err := user.Validate(); err != nil {
		return err
	}
	return ur.db.Save(&user).Error
}

func (ur *UserRepo) DeleteUser(id uint) error {
	user, err := ur.GetUser(id)
	if err != nil {
		return err
	}
	return ur.db.Delete(&user).Error
}
