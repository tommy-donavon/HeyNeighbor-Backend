package utils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/yhung-mea7/HeyNeighbor/account-service/register"
)

type requestOptions struct {
	serviceName string
	methodType  string
	endpoint    string
	register    *register.ConsulClient
	body        []byte
	headers     map[string]string
}

//send new http request
func sendNewRequest(reqOptions *requestOptions) (*http.Response, error) {
	if reqOptions.register == nil {
		return nil, fmt.Errorf("can not look up service with nil consul client")
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
