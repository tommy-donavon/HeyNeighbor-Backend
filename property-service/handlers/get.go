package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
)

func (ph *PropertyHandler) HealthCheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		data.ToJSON(&message{"gud"}, rw)
	}
}

func (ph *PropertyHandler) GetManagerProperties() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		usrCtx, err := ph.ctxHandler.get(r.Context(), "loginInfo")
		if err != nil {
			ph.log.Println("[ERROR] error retrieving login info from context", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		usr, ok := usrCtx.(*data.Tenant)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		props, err := ph.repo.GetAllManagerProperties(usr.Username)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"unable to retrieve any properties for this manager"}, rw)
			return
		}
		data.ToJSON(&props, rw)
	}
}

func (ph *PropertyHandler) GetTenantProperties() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		ph.log.Println("[GET] tenant properties")
		usrCtx, err := ph.ctxHandler.get(r.Context(), "loginInfo")
		if err != nil {
			ph.log.Println("[ERROR] error retrieving login info from context", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		usr, ok := usrCtx.(*data.Tenant)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		props, err := ph.repo.GetAllTenantProperties(usr.Username)
		ph.log.Println(usr.Username)
		if err != nil {
			ph.log.Println("checkpoint")
			ph.log.Println("[ERROR] error adding tenant to property", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"unable to get tenant properties"}, rw)
			return
		}
		data.ToJSON(&props, rw)
	}
}

func (ph *PropertyHandler) GetPropertyByServerCode() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		usrCtx, err := ph.ctxHandler.get(r.Context(), "loginInfo")
		if err != nil {
			ph.log.Println("[ERROR] error retrieving login info from context", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		usr, ok := usrCtx.(*data.Tenant)
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		code := getServerCode(r)
		if len(code) == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"invalid server code"}, rw)
			return
		}
		prop, err := ph.repo.GetPropertyByServerCode(code)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"unable to find property"}, rw)
			return
		}
		for _, t := range prop.Tenants {
			if t.Username == usr.Username {
				data.ToJSON(&prop, rw)
				return
			}
		}
		rw.WriteHeader(http.StatusForbidden)
		data.ToJSON(&message{"you do not belong to this server"}, rw)
	}
}
