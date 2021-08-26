package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
)

// sets Content-type header to application/json for all request
func GlobalContentTypeMiddleware(next http.Handler) http.Handler {
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
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&message{"You are not authorized to make this request"}, rw)
			return
		}
		ctx := ph.ctxHandler.add(r.Context(), "loginInfo", usr)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})
}

// validates that incoming property request data is correct
func (ph *PropertyHandler) ValidatePropertyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prop := data.Property{}
		if err := data.FromJSON(&prop, r.Body); err != nil {
			ph.log.Println("[ERROR] deserializing request body", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to read in request body"}, rw)
			return
		}
		if err := ph.validator.Validate(prop); err != nil {
			ph.log.Println("[ERROR] property is not correctly formated", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{err.Error()}, rw)
			return

		}
		ctx := ph.ctxHandler.add(r.Context(), "propertyInfo", prop)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
