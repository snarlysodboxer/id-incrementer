package main

import (
	"bytes"
	"encoding/json"
	// "fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	// "reflect"
	"strconv"
	"strings"
	"testing"
)

// TODO test paralellism and race conditions, test e2e, test benchmark

type TestID struct {
	ID int
}

type TestError struct {
	Error string
}

func TestListerEndpoint(t *testing.T) {
	// setup
	ids := NewIDMap()
	ids["live"] = map[string]int{"records": 75, "records_other": 67}
	testRouter := ids.SetupRouter()
	request, err := http.NewRequest("GET", "/lister", nil)
	if err != nil {
		t.Error(err)
	}
	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)

	// test for 200 response code
	if response.Code != 200 {
		t.Error("Expected status code 200, got ", response.Code)
	}

	// test for reception of JSON list
	jsonIdMap, err := json.Marshal(ids)
	if err != nil {
		t.Error("Couldn't marshal mocked IDs, this test is broken")
	}
	responseJson := strings.TrimSpace(response.Body.String())
	presetJson := strings.TrimSpace(string(jsonIdMap))
	if responseJson != presetJson {
		t.Errorf("Expected %s, got %s", presetJson, responseJson)
	}
}

func TestGetterEndpoint(t *testing.T) {
	// setup
	ids := NewIDMap()
	testRouter := ids.SetupRouter()
	request, err := http.NewRequest("GET", "/getter/live/records", nil)
	if err != nil {
		t.Error(err)
	}
	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)

	// test for 200 response code
	if response.Code != 200 {
		t.Error("Expected status code 200, got ", response.Code)
	}

	// test for reception of JSON encoded id
	var id TestID
	err = json.Unmarshal(response.Body.Bytes(), &id)
	if err != nil {
		t.Errorf("Unable to unmarshal `%s`", response.Body)
	}
	if id.ID != initialValue {
		t.Errorf("Expected `%d`, got `%d`", initialValue, id.ID)
	}
}

func TestSetterEndpointBadData(t *testing.T) {
	// setup
	number := "56L"
	ids := NewIDMap()
	testRouter := ids.SetupRouter()
	form := url.Values{}
	form.Add("environment", "live")
	form.Add("name", "records_name")
	form.Add("id", number)
	request, err := http.NewRequest("POST", "/setter", bytes.NewBufferString(form.Encode()))
	if err != nil {
		t.Error(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)

	// test for 400 response code
	if response.Code != 400 {
		t.Error("Expected status code 400, got ", response.Code)
	}

	// test for JSON encoded error
	var testError TestError
	err = json.Unmarshal(response.Body.Bytes(), &testError)
	if err != nil {
		t.Errorf("Unable to unmarshal `%s`", response.Body)
	}
	if testError.Error == "" {
		t.Errorf("Expected an `error` field, couldn't find any in `%v`", testError)
	}
}

func TestSetterEndpoint(t *testing.T) {
	// setup
	number := 56
	ids := NewIDMap()
	testRouter := ids.SetupRouter()
	form := url.Values{}
	form.Add("environment", "live")
	form.Add("name", "records_name")
	form.Add("id", strconv.Itoa(number))
	request, err := http.NewRequest("POST", "/setter", bytes.NewBufferString(form.Encode()))
	if err != nil {
		t.Error(err)
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)

	// test for 200 response code
	if response.Code != 200 {
		t.Error("Expected status code 200, got ", response.Code)
	}

	// test for passed JSON encoded id
	var id TestID
	err = json.Unmarshal(response.Body.Bytes(), &id)
	if err != nil {
		t.Errorf("Unable to unmarshal `%s`", response.Body)
	}
	if id.ID != number {
		t.Errorf("Expected `%d`, got `%s`", number, id.ID)
	}

	//// ensure the number remains and gets incremented
	request, err = http.NewRequest("GET", "/getter/live/records_name", nil)
	if err != nil {
		t.Error(err)
	}
	response = httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)

	// test for 200 response code
	if response.Code != 200 {
		t.Error("Expected status code 200, got ", response.Code)
	}

	// test for incremented JSON encoded id
	err = json.Unmarshal(response.Body.Bytes(), &id)
	if err != nil {
		t.Errorf("Unable to unmarshal `%s`", response.Body)
	}
	if id.ID != number+incrementBy {
		t.Errorf("Expected `%d`, got `%s`", number+incrementBy, id.ID)
	}
}

func TestGet(t *testing.T) {
	// setup
	ids := NewIDMap()

	// test for 200 response
	status, id := ids.Get("live", "records")
	if status != 200 {
		t.Error("Expected status code 200, got ", status)
	}

	// test for new initial value
	if id != initialValue {
		t.Errorf("Expected %d, got %d", initialValue, id)
	}

	// test for 200 response
	status, id = ids.Get("live", "records")
	if status != 200 {
		t.Error("Expected status code 200 a second time, got ", status)
	}

	// test for incremented existing value
	if id != initialValue+incrementBy {
		t.Errorf("Expected %d, got %d", initialValue+incrementBy, id)
	}
}

func TestSet(t *testing.T) {
	// setup
	ids := NewIDMap()
	status, id := ids.Set("live", "records", 4242)

	// test for 200 status
	if status != 200 {
		t.Error("Expected status code 200, got ", status)
	}

	// test for expected id
	if id != 4242 {
		t.Error("Expected 4242, got ", id)
	}

	// test for 200 status
	status, id = ids.Set("live", "records", 4242)
	if status != 200 {
		t.Error("Expected status code 200 a second time, got ", status)
	}

	// test for expected id
	if id != 4242 {
		t.Error("Expected 4242 a second time, got ", id)
	}
}
