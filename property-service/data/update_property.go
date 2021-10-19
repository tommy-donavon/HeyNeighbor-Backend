package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	my_json "github.com/yhung-mea7/go-rest-kit/data"
	"go.mongodb.org/mongo-driver/bson"
)

type IPropertyUpdate interface {
	UpdateProperty(string, map[string]interface{}) error
	UpdateTenantInformation(string, map[string]interface{}) error
	PayTenantRent(string, string, float64) (float64, error)
}

// updates base fields for a property
func (pr *PropertyRepo) UpdateProperty(server_code string, updateInfo map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	prop, err := pr.GetPropertyByServerCode(server_code)
	if err != nil {
		return err
	}
	propBytes, err := json.Marshal(prop)
	if err != nil {
		return err
	}
	propMap := map[string]interface{}{}
	err = json.Unmarshal(propBytes, &propMap)
	if err != nil {
		return err
	}

	for key, value := range updateInfo {
		if _, ok := propMap[key]; ok {
			switch key {
			case "property_name":
				v, ok := value.(string)
				if !ok {
					return fmt.Errorf("%s can not be asserted to string", value)
				}
				prop.PropertyName = v
			case "address":
				v, ok := value.(Address)
				if !ok {
					return fmt.Errorf("error asserting value to address")
				}
				prop.PropertyAddress = &v
			case "property_manager":
				v, ok := value.(string)
				if !ok {
					return fmt.Errorf("%s can not be asserted to string", value)
				}
				prop.PropertyManager = v
			case "property_img":
				v, ok := value.(string)
				// v, ok := value.(uint)
				if !ok {
					return fmt.Errorf("%s can not be asserted to string", value)
				}
				prop.PropertyURI = v
			}
		}
	}
	if err := my_json.NewValidator(my_json.ValidationOption{
		Name:      "zip",
		Operation: my_json.NewValidatorFunc(`^\d{5}(-\d{4})?$`),
	}).Validate(prop); err != nil {
		return err
	}

	coll := pr.client.Database(pr.dbName).Collection("properties")
	_, err = coll.UpdateOne(
		ctx,
		bson.M{"_id": prop.ID},
		bson.D{
			{
				Key: "$set",
				Value: bson.D{
					{Key: "property_name", Value: prop.PropertyName},
					{Key: "address", Value: prop.PropertyAddress},
					{Key: "property_manager", Value: prop.PropertyManager},
					// {Key: "num_of_units", Value: prop.NumOfUnits},
				},
			},
		},
	)
	return err
}

//Updates tenant information accross all of their servers
func (pr *PropertyRepo) UpdateTenantInformation(tenantUsername string, updateInfo map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	props, err := pr.GetAllTenantProperties(tenantUsername)
	if err != nil {
		return err
	}
	nilTenant := &Tenant{}
	tenBytes, err := json.Marshal(nilTenant)
	if err != nil {
		return err
	}
	tenMap := map[string]interface{}{}
	if err := json.Unmarshal(tenBytes, &tenMap); err != nil {
		return err
	}
	for _, p := range props {
		for _, ten := range p.Tenants {
			if ten.Username == tenantUsername {
				for key, value := range updateInfo {
					if _, ok := tenMap[key]; ok {
						switch key {
						case "username":
							v, ok := value.(string)
							if !ok {
								return fmt.Errorf("unable to set username")
							}
							ten.Username = v
						case "unit_number":
							v, ok := value.(uint)
							if !ok {
								return fmt.Errorf("unable to set unit number")
							}
							ten.UnitNumber = v
						case "first_name":
							v, ok := value.(string)
							if !ok {
								return fmt.Errorf("unable to set first name")
							}
							ten.FirstName = v
						case "last_name":
							v, ok := value.(string)
							if !ok {
								return fmt.Errorf("unable to set last name")
							}
							ten.LastName = v
						case "email":
							v, ok := value.(string)
							if !ok {
								return fmt.Errorf("unable to set email")
							}
							ten.Email = v
						case "phone_number":
							v, ok := value.(string)
							if !ok {
								return fmt.Errorf("unable to set phone number")
							}
							ten.PhoneNumber = v
						case "profile_uri":
							v, ok := value.(string)
							if !ok {
								return fmt.Errorf("unable to set profile uri")
							}
							ten.ProfileURI = v
						case "user_status":
							v, ok := value.(uint)
							if !ok {
								return fmt.Errorf("unable to set user status")
							}
							ten.UserStatus = v
						}
					}
				}
			}
		}
	}

	for _, p := range props {
		_, err := coll.UpdateOne(
			ctx,
			bson.M{"_id": p.ID},
			bson.D{
				{
					Key:   "$set",
					Value: bson.D{{Key: "tenants", Value: p.Tenants}},
				},
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (pr *PropertyRepo) PayTenantRent(server_code, tenantUsername string, rentPaid float64) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	prop, err := pr.GetPropertyByServerCode(server_code)
	if err != nil {
		return 0, err
	}
	for _, ten := range prop.Tenants {
		if ten.Username == tenantUsername {
			ten.Rent.AmountDue -= rentPaid
			if ten.Rent.AmountDue < 0 {
				return 0, fmt.Errorf("amount provided exceeds rent balance")
			}
			ten.Rent.LastDatePaid = time.Now().String()
			_, err := coll.UpdateOne(
				ctx,
				bson.M{"_id": prop.ID},
				bson.D{
					{
						Key:   "$set",
						Value: bson.D{{Key: "tenants", Value: prop.Tenants}},
					},
				},
			)
			if err != nil {
				return 0, err
			}
			return ten.Rent.AmountDue, nil
		}
	}

	return 0, fmt.Errorf("unable to process rent payment")
}
