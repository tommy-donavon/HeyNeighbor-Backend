package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
	my_json "github.com/yhung-mea7/go-rest-kit/data"
)

func healthcheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		my_json.ToJSON(&message{"gud"}, rw)
	}
}

func getManagerProperties(repo data.IPropertyRead) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		usrCtx, err := instance.ctxHandler.Get(r.Context(), "loginInfo")
		if err != nil {
			instance.log.Println("[ERROR] error retrieving login info from context", err)
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		usr, ok := usrCtx.(*data.Tenant)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		props, err := repo.GetAllManagerProperties(usr.Username)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{"unable to retrieve any properties for this manager"}, rw)
			return
		}
		my_json.ToJSON(&props, rw)
	}
}

func getTenantProperties(repo data.IPropertyRead) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		usrCtx, err := instance.ctxHandler.Get(r.Context(), "loginInfo")
		if err != nil {
			instance.log.Println("[ERROR] error retrieving login info from context", err)
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		usr, ok := usrCtx.(*data.Tenant)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		props, err := repo.GetAllTenantProperties(usr.Username)
		if err != nil {
			instance.log.Println("[ERROR] error adding tenant to property", err)
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{"unable to get tenant properties"}, rw)
			return
		}
		my_json.ToJSON(&props, rw)

	}
}

func getPropertyByServerCode(repo data.IPropertyRead) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {

		usrCtx, err := instance.ctxHandler.Get(r.Context(), "loginInfo")
		if err != nil {
			instance.log.Println("[ERROR] error retrieving login info from context", err)
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		usr, ok := usrCtx.(*data.Tenant)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		code := getServerCode(r)
		if len(code) == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{"invalid server code"}, rw)
			return
		}
		prop, err := repo.GetPropertyByServerCode(code)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{"unable to find property"}, rw)
			return
		}
		for _, t := range prop.Tenants {
			if t.Username == usr.Username {
				my_json.ToJSON(&prop, rw)
				return
			}
		}
		rw.WriteHeader(http.StatusForbidden)
		my_json.ToJSON(&message{"you do not belong to this server"}, rw)
	}
}
