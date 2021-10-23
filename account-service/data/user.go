package data

type (
	User struct {
		Username    string      `json:"username" validate:"required" gorm:"primaryKey"`
		Password    string      `json:"password,omitempty" validate:"required"`
		FirstName   string      `json:"first_name" validate:"required"`
		LastName    string      `json:"last_name" validate:"required"`
		Email       string      `json:"email" validate:"email" gorm:"unique"`
		PhoneNumber string      `json:"phone_number" validate:"phone"`
		ProfileURI  string      `json:"profile_uri"`
		AccountType accountType `json:"account_type" validate:"gte=0,lte=1"`
		UserStatus  status      `json:"user_status" validate:"gte=0,lte=1"`
	}

	accountType uint
	status      uint
)

const (
	ADMIN accountType = iota
	BASE
)
const (
	AVALIABLE status = iota
	UNAVAILABLE
)
