package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

// TODO test paralellism and race conditions, test e2e, test benchmark

func TestListerEndpoint(t *testing.T) {
	// add fake data
	idMap["live"] = map[string]int{"records": 75, "records_other": 67}

	// setup
	testRouter := SetupRouter()
	request, err := http.NewRequest("GET", "/lister", nil)
	if err != nil {
		t.Error(err)
	}
	response := httptest.NewRecorder()

	// test for 200 response code
	testRouter.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Error("Expected status code 200, got ", response.Code)
	}

	// test for reception of JSON list
	jsonIdMap, err := json.Marshal(idMap)
	if err != nil {
		t.Error("Couldn't marshal the fake idMap, your test is broken")
	}
	responseJson := strings.TrimSpace(response.Body.String())
	presetJson := strings.TrimSpace(string(jsonIdMap))
	if responseJson != presetJson {
		t.Errorf("Expected %s, got %s", presetJson, responseJson)
	}
}

func TestGetterEndpoint(t *testing.T) {
	idMap = map[string]map[string]int{}
	testRouter := SetupRouter()
	request, err := http.NewRequest("GET", "/getter/live/records", nil)
	if err != nil {
		t.Error(err)
	}

	response := httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Error("Expected status code 200, got ", response.Code)
	}
	if response.Body.String() != strconv.Itoa(initialValue) {
		t.Errorf("Expected %d, got %s", initialValue, response.Body.String())
	}
}

func TestSetterEndpoint(t *testing.T) {
	number := 56
	testRouter := SetupRouter()
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
	if response.Code != 200 {
		t.Error("Expected status code 200, got ", response.Code)
	}
	if response.Body.String() != strconv.Itoa(number) {
		t.Errorf("Expected %d, got `%s`", number, response.Body.String())
	}

	// ensure the number remains and gets incremented
	request, err = http.NewRequest("GET", "/getter/live/records_name", nil)
	if err != nil {
		t.Error(err)
	}

	response = httptest.NewRecorder()
	testRouter.ServeHTTP(response, request)
	if response.Code != 200 {
		t.Error("Expected status code 200, got ", response.Code)
	}
	if response.Body.String() != strconv.Itoa(number+incrementBy) {
		t.Errorf("Expected %d, got `%s`", number+incrementBy, response.Body.String())
	}
}

func TestGetID(t *testing.T) {
	status, id := getID("live", "records")
	if status != 200 {
		t.Error("Expected status code 200, got ", status)
	}
	if id != strconv.Itoa(initialValue) {
		t.Errorf("Expected %d, got %s", initialValue, id)
	}
	status, id = getID("live", "records")
	if status != 200 {
		t.Error("Expected status code 200 a second time, got ", status)
	}
	if id != strconv.Itoa(initialValue+incrementBy) {
		t.Errorf("Expected %d, got %s", initialValue+incrementBy, id)
	}
}

func TestSetID(t *testing.T) {
	status, id := setID("live", "records", 4242)
	if status != 200 {
		t.Error("Expected status code 200, got ", status)
	}
	if id != "4242" {
		t.Error("Expected 4242, got ", id)
	}
	status, id = setID("live", "records", 4242)
	if status != 200 {
		t.Error("Expected status code 200 a second time, got ", status)
	}
	if id != "4242" {
		t.Error("Expected 4242 a second time, got ", id)
	}
}
