package handlers

import (
	"net/http"
	"strconv"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
)

func (ph *PropertyHandler) CreateProperty() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ph.log.Println("POST Create Property")
		propCtx, err := ph.ctxHandler.get(r.Context(), "propertyInfo")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{err}, rw)
			return
		}
		usrCtx, err := ph.ctxHandler.get(r.Context(), "loginInfo")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{err}, rw)
			return
		}
		prop := propCtx.(data.Property)
		usr := usrCtx.(*data.Tenant)
		if usr.AccountType != 0 {
			rw.WriteHeader(http.StatusForbidden)
			data.ToJSON(&message{"You are not authorized to make this request"}, rw)
			return
		}
		prop.PropertyManager = usr.Username
		prop.ServerCode = ph.repo.GenerateServerCode(prop.PropertyName)
		if err := ph.repo.CreateProperty(&prop); err != nil {
			ph.log.Println("[ERROR] Unable to create property", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to create property"}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}

func (ph *PropertyHandler) AddTenantToProperty() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ph.log.Println("POST ADD TENANT")
		usrCtx, err := ph.ctxHandler.get(r.Context(), "loginInfo")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{err}, rw)
			return
		}
		usr := usrCtx.(*data.Tenant)
		if usr.AccountType != 1 {
			rw.WriteHeader(http.StatusForbidden)
			data.ToJSON(&message{"admins can not be added to property"}, rw)
			return
		}
		rBody := map[string]string{}
		if err := data.FromJSON(&rBody, r.Body); err != nil {
			ph.log.Println("[ERROR] error deseralizing request body", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to process request"}, rw)
			return
		}
		sValue, ok := rBody["server_code"]
		if !ok {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"invalid request body"}, rw)
			return
		}
		uValue, ok := rBody["unit_number"]
		if !ok {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"invalid request body"}, rw)
			return
		}
		unitNumber, err := strconv.ParseUint(uValue, 10, 32)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"unable to process unit number"}, rw)
			return
		}
		rValue, ok := rBody["rent_amount"]
		if !ok {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"invalid request body"}, rw)
			return
		}
		// rentValue, err := strconv.ParseUint(rValue, 10, 32)
		rentValue, err := strconv.ParseFloat(rValue, 32)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"unable to process rent amount"}, rw)
			return
		}
		prop, err := ph.repo.GetPropertyByServerCode(sValue)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"no property found with that code"}, rw)
			return
		}
		usr.UnitNumber = uint(unitNumber)
		usr.Nickname = usr.Username
		usr.Rent.AmountDue = rentValue

		if err := ph.repo.AddTenantToProperty(prop, usr); err != nil {
			ph.log.Println("[ERROR]error adding tenant to property", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"Unable to add tenant to property"}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)

	}
}
