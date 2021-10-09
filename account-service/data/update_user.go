package data

import (
	"encoding/json"
	"fmt"
	"strconv"

	my_json "github.com/yhung-mea7/go-rest-kit/data"
)

type IUserUpdate interface {
	UpdateUser(string, map[string]string) error
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
			case "email":
				user.Email = value
			case "phone_number":
				user.PhoneNumber = value
			case "profile_uri":
				user.ProfileURI = value
			case "user_status":
				iv, err := strconv.ParseUint(value, 10, 32)
				if err != nil {
					return err
				}
				if iv > 3 {
					return fmt.Errorf("%d is not a valid status", iv)
				}
				user.UserStatus = status(iv)
			}
		} else {
			return fmt.Errorf("%s is not a valid field", key)
		}

	}

	if err := my_json.NewValidator(my_json.ValidationOption{
		Name:      "phone",
		Operation: my_json.NewValidatorFunc(`^(\d{1,2}-)?(\d{3}-){2}\d{4}$`),
	}).Validate(user); err != nil {
		return err
	}
	return ur.db.Save(&user).Error
}
