package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
	json "github.com/yhung-mea7/go-rest-kit/data"
)

func createProperty(repo data.IPropertyCreate) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		instance.log.Println("POST Create Property")
		propCtx, err := instance.ctxHandler.Get(r.Context(), "propertyInfo")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{err}, rw)
			return
		}
		usrCtx, err := instance.ctxHandler.Get(r.Context(), "loginInfo")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{err}, rw)
			return
		}
		prop := propCtx.(data.Property)
		usr := usrCtx.(*data.Tenant)
		if usr.AccountType != 0 {
			rw.WriteHeader(http.StatusForbidden)
			json.ToJSON(&message{"You are not authorized to make this request"}, rw)
			return
		}
		prop.PropertyManager = usr.Username
		if err := repo.CreateProperty(&prop); err != nil {
			instance.log.Println("[ERROR] Unable to create property", err)
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{"unable to create property"}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}

func addTenantToProperty(repo data.IPropertyCreate) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		instance.log.Println("POST ADD TENANT")
		usrCtx, err := instance.ctxHandler.Get(r.Context(), "loginInfo")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{err}, rw)
			return
		}
		usr := usrCtx.(*data.Tenant)
		if usr.AccountType != 1 {
			rw.WriteHeader(http.StatusForbidden)
			json.ToJSON(&message{"admins can not be added to property"}, rw)
			return
		}
		rBody := map[string]string{}
		if err := json.FromJSON(&rBody, r.Body); err != nil {
			instance.log.Println("[ERROR] error deseralizing request body", err)
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{"unable to process request"}, rw)
			return
		}
		sValue, ok := rBody["server_code"]
		if !ok {
			rw.WriteHeader(http.StatusBadRequest)
			json.ToJSON(&message{"invalid request body"}, rw)
			return
		}
		// uValue, ok := rBody["unit_number"]
		// if !ok {
		// 	rw.WriteHeader(http.StatusBadRequest)
		// 	json.ToJSON(&message{"invalid request body"}, rw)
		// 	return
		// }
		// unitNumber, err := strconv.ParseUint(uValue, 10, 32)
		// if err != nil {
		// 	rw.WriteHeader(http.StatusBadRequest)
		// 	json.ToJSON(&message{"unable to process unit number"}, rw)
		// 	return
		// }
		// rValue, ok := rBody["rent_amount"]
		// if !ok {
		// 	rw.WriteHeader(http.StatusBadRequest)
		// 	json.ToJSON(&message{"invalid request body"}, rw)
		// 	return
		// }
		// rentValue, err := strconv.ParseFloat(rValue, 32)
		// if err != nil {
		// 	rw.WriteHeader(http.StatusBadRequest)
		// 	json.ToJSON(&message{"unable to process rent amount"}, rw)
		// 	return
		// }
		// usr.UnitNumber = uint(unitNumber)
		usr.Nickname = usr.Username
		// usr.Rent.AmountDue = rentValue
		if err := repo.AddTenantToProperty(sValue, usr); err != nil {
			instance.log.Println("[ERROR]error adding tenant to property", err)
			rw.WriteHeader(http.StatusBadRequest)
			json.ToJSON(&message{"Unable to add tenant to property"}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)

	}
}
