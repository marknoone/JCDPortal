package jcdportal

import (
	"errors"
)

const (
	DefaultHost    = "https://api.jcdecaux.com"
	ParkService    = "vls/v3"
	StationService = "parking/v1"
)

var (
	// NoIDAvailible error is returned when the request executer cannot find
	// enough identification information in the object to carry out a request
	NoIDAvailibleErr = errors.New("no identification field found")

	// NoResultFoundErr is thrown when the specified resource cannot be found
	NoResultFoundErr = errors.New("no result found")

	// TooManyResultsErr is thrown when many results are found when only one is expected
	TooManyResultsErr = errors.New("response carried too many results")
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

func (r apiRequester) Refresh(i interface{}) error { panic("Not Yet Implemented") }

func (r apiRequester) Find(i interface{}, o *RequestOptions) error { panic("Not Yet Implemented") }
