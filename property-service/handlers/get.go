package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
	json "github.com/yhung-mea7/go-rest-kit/data"
)

func (ph *PropertyHandler) HealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		json.ToJSON(&message{"gud"}, rw)
	}
}

func (ph *PropertyHandler) GetManagerProperties() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		usrCtx, err := ph.ctxHandler.Get(r.Context(), "loginInfo")
		if err != nil {
			ph.log.Println("[ERROR] error retrieving login info from context", err)
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		usr, ok := usrCtx.(*data.Tenant)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		props, err := ph.repo.GetAllManagerProperties(usr.Username)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			json.ToJSON(&message{"unable to retrieve any properties for this manager"}, rw)
			return
		}
		json.ToJSON(&props, rw)
	}
}

func (ph *PropertyHandler) GetTenantProperties() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ph.log.Println("[GET] tenant properties")
		usrCtx, err := ph.ctxHandler.Get(r.Context(), "loginInfo")
		if err != nil {
			ph.log.Println("[ERROR] error retrieving login info from context", err)
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		usr, ok := usrCtx.(*data.Tenant)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		props, err := ph.repo.GetAllTenantProperties(usr.Username)
		ph.log.Println(usr.Username)
		if err != nil {
			ph.log.Println("checkpoint")
			ph.log.Println("[ERROR] error adding tenant to property", err)
			rw.WriteHeader(http.StatusBadRequest)
			json.ToJSON(&message{"unable to get tenant properties"}, rw)
			return
		}
		json.ToJSON(&props, rw)
	}
}

func (ph *PropertyHandler) GetPropertyByServerCode() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		usrCtx, err := ph.ctxHandler.Get(r.Context(), "loginInfo")
		if err != nil {
			ph.log.Println("[ERROR] error retrieving login info from context", err)
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		usr, ok := usrCtx.(*data.Tenant)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		code := getServerCode(r)
		if len(code) == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			json.ToJSON(&message{"invalid server code"}, rw)
			return
		}
		prop, err := ph.repo.GetPropertyByServerCode(code)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			json.ToJSON(&message{"unable to find property"}, rw)
			return
		}
		for _, t := range prop.Tenants {
			if t.Username == usr.Username {
				json.ToJSON(&prop, rw)
				return
			}
		}
		rw.WriteHeader(http.StatusForbidden)
		json.ToJSON(&message{"you do not belong to this server"}, rw)
	}
}
