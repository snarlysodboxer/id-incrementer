package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// TODO test paralellism and race conditions, test e2e, test benchmark

func TestGetterEndpoint(t *testing.T) {
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
		t.Errorf("Expected %d, got %s", initialValue, response.Body)
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
