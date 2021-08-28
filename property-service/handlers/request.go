package handlers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yhung-mea7/HeyNeighbor/property-service/data"
	"github.com/yhung-mea7/HeyNeighbor/property-service/register"
)

type requestOptions struct {
	serviceName string
	methodType  string
	endpoint    string
	register    *register.ConsulClient
	body        []byte
	headers     map[string]string
}

// sends GET request to account service to get provided user's information
func (ph *PropertyHandler) getUserInformation(r *http.Request) (*userInformation, error) {
	resp, err := sendNewRequest(&requestOptions{
		serviceName: "account-service",
		methodType:  "GET",
		endpoint:    "",
		register:    ph.register,
		body:        nil,
		headers: map[string]string{
			"Authorization": r.Header.Get("Authorization"),
		},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	userInfo := &userInformation{}
	if err := data.FromJSON(&userInfo, resp.Body); err != nil {
		return nil, err
	}
	return userInfo, nil
}

//send new http request
func sendNewRequest(reqOptions *requestOptions) (*http.Response, error) {
	if reqOptions.register == nil {
		return nil, errors.New("can not look up service with nil consul client")
	}
	ser, err := reqOptions.register.LookUpService(reqOptions.serviceName)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(reqOptions.methodType, ser.GetHTTP()+reqOptions.endpoint, nil)
	if err != nil {
		return nil, err
	}
	if reqOptions.body != nil {

		req.Body = ioutil.NopCloser(strings.NewReader(string(reqOptions.body)))
	}

	for key, value := range reqOptions.headers {
		req.Header.Set(key, value)
	}
	client := &http.Client{}
	return client.Do(req)

}
