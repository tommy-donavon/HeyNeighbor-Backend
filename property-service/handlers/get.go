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

func (ph *PropertyHandler) GetProperties() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		usrCtx, err := ph.ctxHandler.get(r.Context(), "loginInfo")
		if err != nil {
			ph.log.Println("[ERROR] error retrieving login info from context", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to retrieve user information"}, rw)
			return
		}
		usr, ok := usrCtx.(*userInformation)
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

// func (ph *PropertyHandler) GetServerCode() http.HandlerFunc {
// 	return func(rw http.ResponseWriter, r *http.Request) {
// 		addrCTX, err := ph.ctxHandler.get(r.Context(), "addressinfo")
// 		if err != nil {
// 			ph.log.Println("[ERROR] error retreiveing key", err)
// 			data.ToJSON(&message{"error getting address info"}, rw)
// 			return
// 		}
// 		addr, ok := addrCTX.(data.Address)
// 		if !ok {
// 			ph.log.Println("error asserting type data.Address")
// 			data.ToJSON(&message{"unable to assert type to data.Address"}, rw)
// 			return
// 		}
// 		usrCtx, err := ph.ctxHandler.get(r.Context(), "loginInfo")
// 		if err != nil {
// 			rw.WriteHeader(http.StatusInternalServerError)
// 			data.ToJSON(&message{err}, rw)
// 			return
// 		}
// 		usr := usrCtx.(*userInformation)
// 		if usr.UserType != 0 {
// 			rw.WriteHeader(http.StatusForbidden)
// 			data.ToJSON(&message{"You are not authorized to make this request"}, rw)
// 			return
// 		}

// 		prop, err := ph.repo.GetProperty(&addr)
// 		if err != nil {
// 			rw.WriteHeader(http.StatusBadRequest)
// 			data.ToJSON(&message{"property does not exist at that address"}, rw)
// 			return
// 		}
// 		if prop.PropertyManager != usr.Username {
// 			rw.WriteHeader(http.StatusForbidden)
// 			data.ToJSON(&message{"you are not allowed to share this server"}, rw)
// 			return
// 		}
// 		data.ToJSON(&struct {
// 			ServerCode string `json:"server_code"`
// 		}{ServerCode: prop.ServerCode}, rw)
// 	}
// }
