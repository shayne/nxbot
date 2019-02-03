package nxapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	apiFmt = "http://%s/ec2/%s"
)

// API contains information about the API location and caches API data
type API struct {
	ipPort string
	user   string
	pass   string
}

// NewAPI initializes the Nx API given an IP:PORT
func NewAPI(ipPort string, user string, pass string) (*API, error) {
	a := &API{ipPort: ipPort, user: user, pass: pass}
	err := a.testAuth()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *API) testAuth() error {
	req, err := a.newAPIRequest("GET", "getCurrentTime")
	if err != nil {
		return fmt.Errorf("Nx Auth Failure: %v", err)
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Nx Auth Failure: %v", err)
	}
	if resp.StatusCode == 401 {
		return fmt.Errorf("Nx Auth Failure: invalid username or password")
	}
	return nil
}

// GETRequest performs an HTTP GET request to the given endpoint and attempts
// to unmarshal the response to the given data
func (a *API) GETRequest(endpoint string, data interface{}) error {
	req, err := a.newAPIRequest("GET", endpoint)
	if err != nil {
		return fmt.Errorf("GETRequest failed: %v", err)
	}
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("GETRequest failed: %v", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return fmt.Errorf("GETRequest failed: %v", err)
	}
	resp.Body.Close()

	err = json.Unmarshal(buf.Bytes(), &data)
	if err != nil {
		return fmt.Errorf("GETRequest failed: %v", err)
	}
	return nil
}

func (a *API) apiURL(endpoint string) string {
	return fmt.Sprintf(apiFmt, a.ipPort, endpoint)
}

func (a *API) newAPIRequest(method, endpoint string) (*http.Request, error) {
	url := a.apiURL(endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("newAPIRequest failed: %v", err)
	}
	req.SetBasicAuth(a.user, a.pass)
	return req, nil
}
