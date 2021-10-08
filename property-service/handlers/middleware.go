package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
	my_json "github.com/yhung-mea7/go-rest-kit/data"
	request "github.com/yhung-mea7/go-rest-kit/http"
)

// sets Content-type header to application/json for all request
func globalContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

// sends request to account-service to authorize user
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp, err := request.SendNewRequest((&request.RequestOptions{
			ServiceName: "account-service",
			MethodType:  "GET",
			Endpoint:    "",
			Register:    instance.register,
			Body:        nil,
			Headers: map[string]string{
				"Authorization": r.Header.Get("Authorization"),
			},
		}))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{"unable to reach account service"}, rw)
			return
		}
		if resp.StatusCode != http.StatusOK {
			rw.WriteHeader(http.StatusUnauthorized)
			my_json.ToJSON(&message{"you are not authroized to make this request"}, rw)
			return
		}
		defer resp.Body.Close()
		userInfo := &data.Tenant{}
		if err := my_json.FromJSON(&userInfo, resp.Body); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		ctx := instance.ctxHandler.Add(r.Context(), "loginInfo", userInfo)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})
}

// validates that incoming property request data is correct
func validatePropertyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prop := data.Property{}
		if err := my_json.FromJSON(&prop, r.Body); err != nil {
			instance.log.Println("[ERROR] deserializing request body", err)
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{"unable to read in request body"}, rw)
			return
		}
		if err := instance.validator.Validate(prop); err != nil {
			instance.log.Println("[ERROR] property is not correctly formated", err)
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{err.Error()}, rw)
			return

		}
		ctx := instance.ctxHandler.Add(r.Context(), "propertyInfo", prop)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func validateAddressMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		addr := data.Address{}
		if err := my_json.FromJSON(&addr, r.Body); err != nil {
			instance.log.Println("[ERROR] deserializing request body", err)
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{"unable to read in request body"}, rw)
			return
		}
		if err := instance.validator.Validate(addr); err != nil {
			instance.log.Println("[ERROR] property is not correctly formated", err)
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		ctx := instance.ctxHandler.Add(r.Context(), "addressinfo", addr)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}
