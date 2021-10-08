package handlers

import (
	"fmt"
	"net/http"

	models "github.com/yhung-mea7/HeyNeighbor/maintenance-service/data"
	"github.com/yhung-mea7/go-rest-kit/data"
	request "github.com/yhung-mea7/go-rest-kit/http"
)

type (
	userInformation struct {
		Username    string `json:"username" validate:"required"`
		AccountType uint   `json:"account_type" validate:"required"`
		UnitNumber  uint   `json:"unit_number"`
	}

	propInformation struct {
		Tenants         []*userInformation `json:"tenants" validate:"required"`
		PropertyManager string             `json:"property_manager"`
	}
)

func globalContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		resp, err := request.SendNewRequest((&request.RequestOptions{
			ServiceName: "account-service",
			MethodType:  "GET",
			Endpoint:    "",
			Register:    myMainHandler.register,
			Body:        nil,
			Headers: map[string]string{
				"Authorization": r.Header.Get("Authorization"),
			},
		}))
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to reach account service"}, rw)
			return
		}
		if resp.StatusCode != http.StatusOK {
			rw.WriteHeader(http.StatusUnauthorized)
			data.ToJSON(&message{"you are not authroized to make this request"}, rw)
			return
		}
		defer resp.Body.Close()
		userInfo := &userInformation{}
		if err := data.FromJSON(&userInfo, resp.Body); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{err.Error()}, rw)
			return
		}
		if err := myMainHandler.validator.Validate(userInfo); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{err.Error()}, rw)
			return
		}

		ctx := myMainHandler.ctx.Add(r.Context(), "user", userInfo)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)
	})
}

func validateMaintenanceRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		mainRequest := models.MaintenanceRequest{}
		if err := data.FromJSON(&mainRequest, r.Body); err != nil {
			myMainHandler.log.Println("[ERROR] deserializing request body", err)
			rw.WriteHeader(http.StatusInternalServerError)
			data.ToJSON(&message{"unable to read in request body"}, rw)
			return
		}
		if err := myMainHandler.validator.Validate(mainRequest); err != nil {
			myMainHandler.log.Println("[ERROR] property is not correctly formated", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{err.Error()}, rw)
			return
		}
		propertyCode := getServerCode(r)
		propResp, err := request.SendNewRequest(&request.RequestOptions{
			ServiceName: "property-service",
			MethodType:  "GET",
			Endpoint:    fmt.Sprintf("code/%s", propertyCode),
			Register:    myMainHandler.register,
			Body:        nil,
			Headers: map[string]string{
				"Authorization": r.Header.Get("Authorization"),
			},
		})

		if err != nil || propResp.StatusCode != http.StatusOK {
			myMainHandler.log.Panicln("[ERROR] unable to connect to property service", err)
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{"unable to reach property service"}, rw)
			return
		}

		propInfo := &propInformation{}
		if err := data.FromJSON(&propInfo, propResp.Body); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{err.Error()}, rw)
			return
		}
		if err := myMainHandler.validator.Validate(propResp); err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&message{err.Error()}, rw)
			return
		}
		mainRequest.PropertyCode = propertyCode
		ctx := myMainHandler.ctx.Add(r.Context(), "property", propInfo)
		ctx = myMainHandler.ctx.Add(ctx, "request", mainRequest)
		r = r.WithContext(ctx)
		next.ServeHTTP(rw, r)

	})
}

// func updateValidationMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		uCode, err := getRequestId(r)
// 		if err != nil {
// 			rw.WriteHeader(http.StatusBadRequest)
// 			data.ToJSON(&message{err.Error()}, rw)
// 			return
// 		}
// 		userCtx, err := myMainHandler.ctx.Get(r.Context(), "user")
// 		if err != nil {
// 			rw.WriteHeader(http.StatusInternalServerError)
// 			data.ToJSON(&message{err.Error()}, rw)
// 			return
// 		}
// 		usr := userCtx.(*userInformation)
// 		request, err := myMainHandler.repo.GetMaintenanceRequest(uCode)
// 		if err != nil {
// 			rw.WriteHeader(http.StatusBadRequest)
// 			data.ToJSON(&message{err.Error()}, rw)
// 			return
// 		}
// 		if (usr.AccountType == 0 && request.Admin != usr.Username) || (usr.AccountType == 1 && request.Tenant != usr.Username) {
// 			rw.WriteHeader(http.StatusForbidden)
// 			data.ToJSON(&message{"you are not allowed to alter this service request"}, rw)
// 			return
// 		}
// 		next.ServeHTTP(rw, r)
// 	})
// }
