package jcdportal

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	DefaultHost = "https://api.jcdecaux.com"
	StdService  = "vls/v3"
	ParkService = "parking/v1"
)

var (
	// UnauthorizedErr reurned when API key cannot access resource
	UnauthorizedErr = errors.New("request error: apikey unauthorized")

	// NoResourceFoundErr is thrown when the specified resource cannot be found
	NoResourceFoundErr = errors.New("request error: no resource found")
)

type apiRequester struct {
	key  string
	host string
}

func NewRequester(key string) *apiRequester { return &apiRequester{key: key} }
func (r *apiRequester) WithHost(h string) {
	r.host = h
}

func (r *apiRequester) GetContracts() (*[]Contract, error) {
	var (
		// value considers current contract amount of 27 total
		result = make([]Contract, 0, 50)
		uri    = fmt.Sprintf("%s/%s/contracts", r.host, StdService)
	)

	err := r.makeRequest(uri, &result)
	return &result, err
}

func (r *apiRequester) GetStations() (*[]Station, error) {
	var (
		result = make([]Station, 0, 3000) // value considers current 2620 total
		uri    = fmt.Sprintf("%s/%s/stations", r.host, StdService)
	)

	err := r.makeRequest(uri, &result)
	return &result, err
}

func (r *apiRequester) GetParks(contractName string) (*[]Park, error) {
	var (
		result = make([]Park, 0, 250)
		uri    = fmt.Sprintf("%s/%s/contracts/%s/parks", r.host, ParkService, contractName)
	)

	err := r.makeRequest(uri, &result)
	return &result, err
}

func (r *apiRequester) GetStationsInContract(contractName string) (*[]Station, error) {

	var (
		// value considers current largest contract in Lyon with 400 total
		result = make([]Station, 0, 500)
		uri    = fmt.Sprintf("%s/%s/stations?contract=%s", r.host, StdService, contractName)
	)

	err := r.makeRequest(uri, &result)
	return &result, err
}

func (r *apiRequester) GetPark(contractName string, number int) (*Park, error) {
	var (
		result Park
		uri    = fmt.Sprintf("%s/%s/contracts/%s/parks/%d", r.host, ParkService, contractName, number)
	)

	err := r.makeRequest(uri, &result)
	return &result, err
}

func (r *apiRequester) GetStation(contractName string, number int) (*Station, error) {
	var (
		result Station
		uri    = fmt.Sprintf("%s/%s/stations/%d?contract=%s", r.host, StdService, number, contractName)
	)

	err := r.makeRequest(uri, &result)
	return &result, err
}

func (r *apiRequester) makeRequest(uri string, dest interface{}) error {
	if r.host == "" {
		r.host = DefaultHost
	}

	reqURL, err := url.Parse(uri)
	if err != nil {
		return err
	}

	qv := reqURL.Query()
	qv.Set("apiKey", r.key)
	reqURL.RawQuery = qv.Encode()

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200:
		return json.NewDecoder(resp.Body).Decode(dest)
	case 403:
		return UnauthorizedErr
	case 404:
		return NoResourceFoundErr
	default:
		return fmt.Errorf(
			"request error occured: recieved http status code (%d)",
			resp.StatusCode)
	}
}
