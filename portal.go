package jcdportal

import (
	"encoding/json"
	"fmt"
	"reflect"
	"errors"

	"net/http"
	"net/url"
)

const (
	DefaultHost = "https://api.jcdecaux.com"
	StdService  = "vls/v3"
)

var (
	// NoIDAvailible error is returned when the request executer cannot find
	// enough identification information in the object to carry out a request
	NoIDAvailibleErr = errors.New("no identification field found")

	// NoResultFoundErr is thrown when the specified resource cannot be found
	NoResultFoundErr = errors.New("no result found")

	// TooManyResultsErr is thrown when many results are found when only one
	// is expected
	TooManyResultsErr = errors.New("response carried too many results")

	// ExpectedPointerErr is thrown when a non-pointer value is passed
	// into the Find/Refresh functions
	ExpectedPointerErr = errors.New("expected pointer variable")

	// UnrecognisedTypeErr is thrown when the consumer passes a non JCD portal
	// data type into the Find/Refresh functions
	UnrecognisedTypeErr = errors.New(
		"unexpected type: function only supports portal data types")
)

type apiRequester struct {
	key  string
	host string
}

type RequestOptions struct {
	Number       int
	ContractName string
}

func NewRequester(key string) apiRequester { return apiRequester{key: key} }
func (r apiRequester) WithHost(h string) apiRequester {
	r.host = h
	return r
}

func (r apiRequester) Refresh(i interface{}) error { panic("Not yet implemented") }
func (r apiRequester) Find(i interface{}, o *RequestOptions) error {
	if reflect.ValueOf(i).Kind() != reflect.Ptr {
		return ExpectedPointerErr
	}

	reqURL, err := r.determineUrl(i, o)
	if err != nil {
		return err
	}

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

	return json.NewDecoder(resp.Body).Decode(i)
}

// determineUrl creates a URL object based off of the API spec
// found here: https://developer.jcdecaux.com/#/opendata/vls?page=dynamic
func (r apiRequester) determineUrl(i interface{}, o *RequestOptions) (*url.URL, error) {
	if r.host == "" {
		r.host = DefaultHost
	}

	var (
		baseType reflect.Type
		reqUrl   *url.URL
		qv       url.Values
	)

	iElem := reflect.ValueOf(i).Elem()
	if iElem.Kind() == reflect.Slice || iElem.Kind() == reflect.Array {
		baseType = reflect.ValueOf(i).Elem().Index(0).Type()
	} else {
		baseType = reflect.ValueOf(i).Elem().Type()
	}

	reqUrl, err := url.Parse(fmt.Sprintf("%s/%s", r.host, StdService))
	if err != nil {
		return nil, err
	}

	qv, err = url.ParseQuery(fmt.Sprintf("apiKey=%S", r.key))
	if err != nil {
		return nil, err
	}

	switch baseType {
	case reflect.TypeOf(Contract{}):
		reqUrl.RawQuery = qv.Encode()
		return reqUrl.Parse("/contracts")
	case reflect.TypeOf(Station{}):
		reqUrl.Path = reqUrl.Path + "/stations"
	default:
		return nil, UnrecognisedTypeErr
	}

	if o.Number != 0 {
		reqUrl.Path = reqUrl.Path + fmt.Sprintf("/%d", o.Number)
	}

	if o.ContractName != "" {
		qv.Add("contract", o.ContractName)
	}

	reqUrl.RawQuery = qv.Encode()
	return reqUrl.Parse("")
}
