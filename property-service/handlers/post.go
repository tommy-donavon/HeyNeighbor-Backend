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
		usr := usrCtx.(userInformation)
		if usr.UserType != 0 {
			rw.WriteHeader(http.StatusForbidden)
			data.ToJSON(&message{"You are not authorized to make this request"}, rw)
			return
		}
		prop.PropertyManager = usr.Username
		if err := ph.repo.CreateProperty(&prop); err != nil {
			ph.log.Println("[ERROR] Unable to create property", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to create property"}, rw)
			return
		}
		rw.WriteHeader(http.StatusNoContent)
	}
}
