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
	testContractName = "testing"
	testNumber       = 84
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

		fmt.Println(key)
		if res, ok := rc[key]; ok {
			w.Write(res)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Something bad happened!"))
		}
	}))
	defer srv.Close()

	// Test Cases
	r := jcd.NewRequester("123456").WithHost(srv.URL)
	tcs := []struct {
		target   interface{}
		expected interface{}
		opts     *jcd.RequestOptions
	}{
		{target: &[]jcd.Contract{}, opts: nil, expected: bytesWithNoErr(json.Marshal([]jcd.Contract{jcd.DummyContract}))},
		{target: &[]jcd.Station{}, opts: nil, expected: bytesWithNoErr(json.Marshal([]jcd.Station{jcd.DummyStation}))},
		{target: &[]jcd.Station{}, opts: &jcd.RequestOptions{ContractName: testContractName}, expected: bytesWithNoErr(json.Marshal([]jcd.Station{jcd.DummyStation}))},
		{target: &[]jcd.Park{}, opts: &jcd.RequestOptions{ContractName: testContractName}, expected: bytesWithNoErr(json.Marshal([]jcd.Park{jcd.DummyPark}))},
		{target: &jcd.Station{}, opts: &jcd.RequestOptions{ContractName: testContractName, Number: testNumber}, expected: bytesWithNoErr(json.Marshal(jcd.DummyStation))},
		{target: &jcd.Park{}, opts: &jcd.RequestOptions{ContractName: testContractName, Number: testNumber}, expected: bytesWithNoErr(json.Marshal(jcd.DummyPark))},
	}

	for _, tc := range tcs {
		err := r.Find(tc.target, tc.opts)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(tc.target, tc.expected) {
			t.Fatalf("Got result %+v\nExpected %+v", tc.target, tc.expected)
		}
	}
}

func TestRefresh(t *testing.T) {

	// Server
	hasRequested := false
	modifiedDummyStation := jcd.DummyStation
	modifiedDummyStation.Name = "Dummy"

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

		m := jcd.DummyStation
		if hasRequested {
			m = modifiedDummyStation
		}

		res, err := json.Marshal(m)
		if err != nil {
			t.Fatal(err)
		}
		w.Write(res)
	}))
	defer srv.Close()

	var s jcd.Station
	r := jcd.NewRequester(apiTestKey).WithHost(srv.URL)

	err := r.Find(s, &jcd.RequestOptions{Number: testNumber, ContractName: testContractName})
	switch {
	case err != nil:
		t.Fatal(err)
	case !reflect.DeepEqual(s, jcd.DummyStation):
		t.Fatalf("Got result %+v\nExpected %+v", s, jcd.DummyStation)
	}

	err = r.Refresh(s)
	switch {
	case err != nil:
		t.Fatal(err)
	case !reflect.DeepEqual(s, modifiedDummyStation):
		t.Fatalf("Got result %+v\nExpected %+v", s, jcd.DummyStation)
	}

}
