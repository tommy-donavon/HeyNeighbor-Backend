package handlers

import (
	"net/http"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
	json "github.com/yhung-mea7/go-rest-kit/data"
	request "github.com/yhung-mea7/go-rest-kit/http"
)

// sends GET request to account service to get provided user's information
func (ph *PropertyHandler) getUserInformation(r *http.Request) (*data.Tenant, error) {
	resp, err := request.SendNewRequest((&request.RequestOptions{
		ServiceName: "account-service",
		MethodType:  "GET",
		Endpoint:    "",
		Register:    ph.register,
		Body:        nil,
		Headers: map[string]string{
			"Authorization": r.Header.Get("Authorization"),
		},
	}))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	userInfo := &data.Tenant{}
	if err := json.FromJSON(&userInfo, resp.Body); err != nil {
		return nil, err
	}
	return userInfo, nil
}
