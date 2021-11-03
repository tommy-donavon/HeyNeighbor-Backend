package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/maintenance-service/data"
	my_json "github.com/yhung-mea7/go-rest-kit/data"
)

func healthcheck() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		my_json.ToJSON(&message{"service healthy"}, rw)
	}
}

func getPropertyMaintenanceRequest(repo data.IMaintenanceRead) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		serverCode := getServerCode(r)
		userCtx, err := myMainHandler.ctx.Get(r.Context(), "user")
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		usr := userCtx.(*userInformation)
		requests, err := repo.GetPropertyMaintenanceRequest(serverCode)
		if err != nil {
			myMainHandler.log.Println("[ERROR] fetching property maintenance requests", err)
			rw.WriteHeader(http.StatusBadRequest)
			my_json.ToJSON(&message{err.Error()}, rw)
			return
		}
		accountType := uint(1)
		if usr.AccountType == &accountType {
			filteredRequest := []*data.MaintenanceRequest{}
			for _, r := range requests {
				if r.Tenant == usr.Username {
					filteredRequest = append(filteredRequest, r)
				}
			}
			requests = filteredRequest
		}
		my_json.ToJSON(&requests, rw)
	}
}
