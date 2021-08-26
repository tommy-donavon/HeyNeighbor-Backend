package data

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type (
	User struct {
		Username    string      `json:"username" validate:"required" gorm:"primaryKey"`
		Password    string      `json:"password,omitempty" validate:"required"`
		FirstName   string      `json:"first_name" validate:"required"`
		LastName    string      `json:"last_name" validate:"required"`
		UnitNumber  uint        `json:"unit_number" validate:"required"`
		Email       string      `json:"email" validate:"email" gorm:"unique"`
		PhoneNumber string      `json:"phone_number" validate:"phone" gorm:"unique"`
		ProfileURI  string      `json:"profile_uri"`
		AccountType accountType `json:"account_type" validate:"gte=0,lte=2"`
		UserStatus  status      `json:"user_status" validate:"gte=0,lte=3"`
	}

	accountType uint
	status      uint
)

const (
	ADMIN accountType = iota
	BASE
	SUB
)
const (
	ONLINE status = iota
	OFFLINE
	DONOTDISTURB
	IDLE
)

func (ur *UserRepo) CreateUser(user *User) error {
	hash, err := hashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash
	return ur.db.Create(user).Error
}

func (ur *UserRepo) GetUser(username string) (*User, error) {
	user := User{}
	err := ur.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (ur *UserRepo) UpdateUser(username string, updateInfo map[string]string) error {
	user, err := ur.GetUser(username)
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

func (ur *UserRepo) DeleteUser(username string) error {
	user, err := ur.GetUser(username)
	if err != nil {
		return err
	}
	return ur.db.Delete(&user).Error
}
