package jcdportal_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	jcd "github.com/marknoone/JCDPortal"
)

var (
	apiTestKey       = "123456"
	testContractName = "Lyon"
	testNumber       = 123
)

func TestPortal(t *testing.T) {

	bytesWithNoErr := func(b []byte, err error) []byte {
		if err != nil {
			t.Fatal(err)
		}
		return b
	}

	// Response Cases
	rc := map[string][]byte{
		fmt.Sprintf("/vls/v3/contracts"):                                               bytesWithNoErr(json.Marshal([]jcd.Contract{jcd.DummyContract})),
		fmt.Sprintf("/vls/v3/stations"):                                                bytesWithNoErr(json.Marshal([]jcd.Station{jcd.DummyStation})),
		fmt.Sprintf("/vls/v3/stations?contract=%s", testContractName):                  bytesWithNoErr(json.Marshal([]jcd.Station{jcd.DummyStation})),
		fmt.Sprintf("/parking/v1/contracts/%s/parks", testContractName):                bytesWithNoErr(json.Marshal([]jcd.Park{jcd.DummyPark})),
		fmt.Sprintf("/vls/v3/stations/%d?contract=%s", testNumber, testContractName):   bytesWithNoErr(json.Marshal(jcd.DummyStation)),
		fmt.Sprintf("/parking/v1/contracts/%s/parks/%d", testContractName, testNumber): bytesWithNoErr(json.Marshal(jcd.DummyPark)),
	}

	// Server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qv := r.URL.Query()
		if qv["apiKey"][0] != apiTestKey {
			w.Write([]byte("API Key Not Recognised"))
		}

		key := r.URL.Path
		if c, hasContract := qv["contract"]; hasContract {
			key = fmt.Sprintf("%s?contract=%s", key, c[0])
		}

		if res, ok := rc[key]; ok {
			w.Write(res)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
		}
	}))
	defer srv.Close()

	r := jcd.NewRequester(apiTestKey)
	r.WithHost(srv.URL)
	wrapErr := func(i interface{}, err error) interface{} {
		if err != nil {
			t.Fatal(err)
		}
		return i
	}

	// Test Cases
	tcs := []struct {
		name     string
		target   interface{}
		expected interface{}
	}{
		{
			name:     "Test case: contracts",
			target:   wrapErr(r.GetContracts()),
			expected: &[]jcd.Contract{jcd.DummyContract},
		},
		{
			name:     "Test case: stations",
			target:   wrapErr(r.GetStations()),
			expected: &[]jcd.Station{jcd.DummyStation},
		},
		{
			name:     "Test case: stations of contract",
			target:   wrapErr(r.GetStationsInContract(testContractName)),
			expected: &[]jcd.Station{jcd.DummyStation},
		},
		{
			name:     "Test case: parks of contract",
			target:   wrapErr(r.GetParks(testContractName)),
			expected: &[]jcd.Park{jcd.DummyPark},
		},
		{
			name:     "Test case: station",
			target:   wrapErr(r.GetStation(testContractName, testNumber)),
			expected: &jcd.DummyStation,
		},
		{
			name:     "Test case: park",
			target:   wrapErr(r.GetPark(testContractName, testNumber)),
			expected: &jcd.DummyPark,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if !reflect.DeepEqual(tc.target, tc.expected) {
				t.Fatalf("Got result %+v\nExpected %+v", tc.target, tc.expected)
			}
		})
	}
}

func TestRefresh(t *testing.T) {

	// Server
	hasRequested := false
	modifiedDummyStation := jcd.DummyStation

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		qv := r.URL.Query()
		key := r.URL.Path
		if c, hasContract := qv["contract"]; hasContract {
			key = fmt.Sprintf("%s?contract=%s", key, c[0])
		} else {
			t.Fatal("No contract specified")
		}

		if key != fmt.Sprintf("/vls/v3/stations/%d?contract=%s", testNumber, testContractName) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"Resource not found"}`))
		}

		m := modifiedDummyStation
		if hasRequested {
			m.Name = "Dummy"
		}

		res, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		w.Write(res)
	}))
	defer srv.Close()

	r := jcd.NewRequester(apiTestKey)
	r.WithHost(srv.URL)
	s, err := r.GetStation(testContractName, testNumber)
	switch {
	case err != nil:
		t.Fatal(err)
	case !reflect.DeepEqual(s, &jcd.DummyStation):
		t.Fatalf("Got result %+v\nExpected %+v", s, jcd.DummyStation)
	}

	err = s.Refresh(r)
	switch {
	case err != nil:
		t.Fatal(err)
	case !reflect.DeepEqual(s, &modifiedDummyStation):
		t.Fatalf("Got result %+v\nExpected %+v", s, jcd.DummyStation)
	}

}
