package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type serviceMock struct {
	returnValue *DeviceWithId
	returnError error
}

func (s *serviceMock) Create(device *Device) (*DeviceWithId, error) {
	return s.returnValue, s.returnError
}

func (s *serviceMock) Get(id int) (*DeviceWithId, error) {
	return s.returnValue, s.returnError
}

func TestCreateDevice_Pass(t *testing.T) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	inputJson := `{"name": "name", "interval": 5, "value": 1.4}`
	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandlePost)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusCreated, rr.Code)
	require.JSONEq(t, `{"id":1, "interval":5, "name":"name", "value":1.4}`, rr.Body.String())
}

func TestCreateDevice_InvalidJson(t *testing.T) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	inputJson := `{"name": 1}`
	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandlePost)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"message":"invalid json body","error":"bad_request"}`, rr.Body.String())
}

func TestCreateDevice_NotValidated(t *testing.T) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	inputJson := `{"name": "test name", "interval": -1}`
	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandlePost)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"message":"Request cannot be validate: interval has to be greater than 0;","error":"bad_request"}`, rr.Body.String())
}

func TestCreateDevice_InternalServerError(t *testing.T) {
	deviceService := &serviceMock{returnError: NewInternalServerError("database error")}
	deviceController := DeviceController{deviceService}

	inputJson := `{"name": "name", "interval": 5, "value": 1.4}`
	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandlePost)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"message":"database error","error":"server_error"}`, rr.Body.String())
}

//
func TestGetDevice_Pass(t *testing.T) {
	repo := NewInMemRepo()
	repo.Save(&Device{Name: "name", Interval: 5, Value: 1.4})
	deviceController := DeviceController{&DeviceService{repo}}
	req, err := http.NewRequest(http.MethodGet, "/devices", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusOK, rr.Code)
	require.JSONEq(t, `{"id":1, "interval":5, "name":"name", "value":1.4}`, rr.Body.String())
}

func TestGetDevice_BadRequest(t *testing.T) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	req, err := http.NewRequest(http.MethodGet, "/devices", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "string",
	})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"bad_request", "message":"id must be a string"}`, rr.Body.String())
}

func TestGetDevice_NotFound(t *testing.T) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	req, err := http.NewRequest(http.MethodGet, "/devices", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusNotFound, rr.Code)
	require.JSONEq(t, `{"error":"not_found", "message":"device not found"}`, rr.Body.String())
}

func TestGetDevice_InternalServerError(t *testing.T) {
	deviceService := &serviceMock{returnError: NewInternalServerError("database error")}
	deviceController := DeviceController{deviceService}

	inputJson := `{"name": "name", "interval": 5, "value": 1.4}`
	req, err := http.NewRequest(http.MethodGet, "/devices/1", bytes.NewBufferString(inputJson))
	req = mux.SetURLVars(req, map[string]string{
		"id": "1",
	})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"message":"database error","error":"server_error"}`, rr.Body.String())
}
