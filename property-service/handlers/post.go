package handlers

import (
	"net/http"

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
		usr := usrCtx.(*userInformation)
		if usr.UserType != 0 {
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

//TODO finish validating server code
func (ph *PropertyHandler) AddTenantToProperty() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ph.log.Println("POST ADD TENANT")
		usrCtx, err := ph.ctxHandler.get(r.Context(), "loginInfo")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{err}, rw)
			return
		}
		usr := usrCtx.(*userInformation)
		if usr.UserType != 1 {
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
		rValue, ok := rBody["server_code"]
		if !ok {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"invalid request body"}, rw)
			return
		}
		prop, err := ph.repo.GetPropertyByServerCode(rValue)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"no property found with that code"}, rw)
			return
		}
		tenantInfo := &data.Tenant{
			Username:   usr.Username,
			Nickname:   usr.Username,
			UnitNumber: 0,
			ProfileURI: usr.ProfileUri,
		}
		if err := ph.repo.AddTenantToProperty(prop, tenantInfo); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{err}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)

	}
}
