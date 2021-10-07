package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
	json "github.com/yhung-mea7/go-rest-kit/data"
)

// sets Content-type header to application/json for all request
func (ph *PropertyHandler) GlobalContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

// sends request to account-service to authorize user
func (ph *PropertyHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		usr, err := ph.getUserInformation(r)
		if err != nil {
			ph.log.Println("[ERROR] error getting user information", err)
			rw.WriteHeader(http.StatusUnauthorized)
			json.ToJSON(&message{"You are not authorized to make this request"}, rw)
			return
		}
		ctx := ph.ctxHandler.Add(r.Context(), "loginInfo", usr)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})
}

// validates that incoming property request data is correct
func (ph *PropertyHandler) ValidatePropertyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prop := data.Property{}
		if err := json.FromJSON(&prop, r.Body); err != nil {
			ph.log.Println("[ERROR] deserializing request body", err)
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{"unable to read in request body"}, rw)
			return
		}
		if err := ph.validator.Validate(prop); err != nil {
			ph.log.Println("[ERROR] property is not correctly formated", err)
			rw.WriteHeader(http.StatusBadRequest)
			json.ToJSON(&message{err.Error()}, rw)
			return

		}
		ctx := ph.ctxHandler.Add(r.Context(), "propertyInfo", prop)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func (ph *PropertyHandler) ValidateAddressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		addr := data.Address{}
		if err := json.FromJSON(&addr, r.Body); err != nil {
			ph.log.Println("[ERROR] deserializing request body", err)
			rw.WriteHeader(http.StatusInternalServerError)
			json.ToJSON(&message{"unable to read in request body"}, rw)
			return
		}
		if err := ph.validator.Validate(addr); err != nil {
			ph.log.Println("[ERROR] property is not correctly formated", err)
			rw.WriteHeader(http.StatusBadRequest)
			json.ToJSON(&message{err.Error()}, rw)
			return
		}
		ctx := ph.ctxHandler.Add(r.Context(), "addressinfo", addr)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
