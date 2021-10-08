package data

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Property struct {
		ID              primitive.ObjectID `bson:"_id,omitempty" json:"-"`
		PropertyName    string             `bson:"property_name" json:"property_name" validate:"required"`
		PropertyAddress *Address           `bson:"address" json:"address" validate:"required,dive,required"`
		PropertyManager string             `bson:"property_manager" json:"property_manager"`
		NumOfUnits      uint               `bson:"num_of_units" json:"num_of_units" validate:"required"`
		ServerCode      string             `bson:"server_code" json:"server_code"`
		Tenants         []*Tenant          `bson:"tenants" json:"tenants" validate:"dive"`
		Channels        []string           `bson:"channels" json:"channels"`
	}

	Tenant struct {
		Username    string        `json:"username"`
		Nickname    string        `json:"nickname"`
		UnitNumber  uint          `json:"unit_number"`
		FirstName   string        `json:"first_name"`
		LastName    string        `json:"last_name"`
		Email       string        `json:"email"`
		PhoneNumber string        `json:"phone_number"`
		ProfileURI  string        `json:"profile_uri"`
		AccountType uint          `json:"account_type"`
		UserStatus  uint          `json:"user_status"`
		Rent        RentAgreement `json:"rent_agreement"`
	}

	RentAgreement struct {
		AmountDue    float64 `json:"amount_due"`
		LastDatePaid string  `json:"last_date_paid"`
		// AmountMonthly float64 `json:"amount_monthly"`
	}
	Address struct {
		StreetAddress string `json:"street_address" validate:"required"`
		City          string `json:"city" validate:"required"`
		State         string `json:"state" validate:"required"`
		ZipCode       string `json:"zip_code" validate:"required,zip"`
	}
)

const (
	//general chat
	general string = "General"
	//landlord announcements
	announcements string = "Announcements"
	//events
	events string = "Events"
)

//TODO add go routine to update user rent payments
