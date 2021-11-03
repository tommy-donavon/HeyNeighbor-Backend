package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/maintenance-service/data"
	my_json "github.com/yhung-mea7/go-rest-kit/data"
)

func createMaintenanceRequest(repo data.IMaintenanceCreate) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		propCtx, err := myMainHandler.ctx.Get(r.Context(), "property")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		userCtx, err := myMainHandler.ctx.Get(r.Context(), "user")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		requestCtx, err := myMainHandler.ctx.Get(r.Context(), "request")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		prop := propCtx.(*propInformation)
		usr := userCtx.(*userInformation)
		request := requestCtx.(data.MaintenanceRequest)

		accountType := uint(1)
		for _, t := range prop.Tenants {
			if t.Username == usr.Username && usr.AccountType == &accountType {
				request.Tenant = usr.Username
				request.UnitNumber = usr.UnitNumber
				request.Admin = prop.PropertyManager
				if err := repo.CreateMaintenanceRequest(&request); err != nil {
					rw.WriteHeader(http.StatusInternalServerError)
					my_json.ToJSON(&message{err.Error()}, rw)
					return
				}
				rw.WriteHeader(http.StatusNoContent)
				return
			}
		}
		rw.WriteHeader(http.StatusForbidden)
		my_json.ToJSON(&message{"you are not able to create a request for this property"}, rw)
	}
}
