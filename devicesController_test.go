package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type serviceMock struct {
	returnValue *DeviceWithId
	returnError *MessageErr
}

func (s *serviceMock) Create(device *Device) (*DeviceWithId, *MessageErr) {
	return s.returnValue, s.returnError
}

func TestCreateMessage_Pass(t *testing.T) {
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

func TestCreateMessage_InvalidJson(t *testing.T) {
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

func TestCreateMessage_NotValidated(t *testing.T) {
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

func TestCreateMessage_InternalServerError(t *testing.T) {
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
