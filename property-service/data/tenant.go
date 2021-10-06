package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Tenant struct {
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

type RentAgreement struct {
	AmountDue    float64 `json:"amount_due"`
	LastDatePaid string  `json:"last_date_paid"`
	// AmountMonthly float64 `json:"amount_monthly"`
}

//TODO add method to update tenant server nickname

// add tenant
func (pr *PropertyRepo) AddTenantToProperty(prop *Property, ten *Tenant) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	for _, t := range prop.Tenants {
		if t.Username == ten.Username {
			return fmt.Errorf("%s is already in server", ten.Username)
		} else if t.Nickname == ten.Nickname {
			return fmt.Errorf("%s is already taken on server", ten.Nickname)
		}
	}
	prop.Tenants = append(prop.Tenants, ten)
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
	return err
}

func (pr *PropertyRepo) GetAllTenantProperties(tenantUsername string) ([]*Property, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	result, err := coll.Find(
		ctx,
		bson.M{"tenants.username": tenantUsername},
	)
	if err != nil {
		return nil, err
	}
	results := []*Property{}
	for result.Next(ctx) {
		prop := Property{}
		if err := result.Decode(&prop); err != nil {
			return nil, err
		}
		results = append(results, &prop)
	}
	return results, nil
}

func (pr *PropertyRepo) RemoveTenantFromProperty(serverCode, tenUsername string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	coll := pr.client.Database(pr.dbName).Collection("properties")
	prop, err := pr.GetPropertyByServerCode(serverCode)
	if err != nil {
		return err
	}
	inList := false
	for i, t := range prop.Tenants {
		if t.Username == tenUsername {
			prop.Tenants[i] = prop.Tenants[len(prop.Tenants)-1]
			prop.Tenants = prop.Tenants[:len(prop.Tenants)-1]
			inList = true
		}
	}
	if !inList {
		return fmt.Errorf("%s is not a tenant of this property", tenUsername)
	}
	_, err = coll.UpdateOne(
		ctx,
		bson.M{"_id": prop.ID},
		bson.D{
			{
				Key:   "$set",
				Value: bson.D{{Key: "tenants", Value: prop.Tenants}},
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

// func(pr *PropertyRepo) SetTenantUserName(server_code, tenantUsername, newTenantNickname string) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
// 	defer cancel()
// 	coll := pr.client.Database(pr.dbName).Collection("properties")
// 	prop, err := pr.GetPropertyByServerCode(server_code)
// 	if err != nil {
// 		return err
// 	}

// }

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
