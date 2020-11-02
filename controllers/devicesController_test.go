package controllers

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"tsm/models"
	"tsm/services"
)

var (
	createDevice func(message *models.Device) (*models.DeviceWithId, models.MessageErr)
)

type deviceMock struct{}

func (dm *deviceMock) CreateDevice(message *models.Device) (*models.DeviceWithId, models.MessageErr) {
	return createDevice(message)
}

func TestCreateMessage_Pass(t *testing.T) {
	services.DeviceService = &deviceMock{}
	createDevice = func(message *models.Device) (*models.DeviceWithId, models.MessageErr) {
		return &models.DeviceWithId{
			Id: 1,
			Device: &models.Device{
				Name:     "name",
				Interval: 5,
				Value:    1.4,
			},
		}, nil
	}
	inputJson := `{"name": "name", "interval": 5, "value": 1.4}`
	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlePost)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusCreated, rr.Code)
	require.JSONEq(t, `{"id":1, "interval":5, "name":"name", "value":1.4}`, rr.Body.String())
}

func TestCreateMessage_InvalidJson(t *testing.T) {
	inputJson := `{"name": 1}`
	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlePost)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"message":"invalid json body","status":400,"error":"bad_request"}`, rr.Body.String())
}

func TestCreateMessage_NotValidated(t *testing.T) {
	services.DeviceService = &deviceMock{}
	createDevice = func(message *models.Device) (*models.DeviceWithId, models.MessageErr) {
		return &models.DeviceWithId{}, models.NewBadRequestError("Request cannot be validate: interval has to be greater than 0;")
	}
	inputJson := `{"name": "test name", "interval": -1}`
	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlePost)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"message":"Request cannot be validate: interval has to be greater than 0;","status":400,"error":"bad_request"}`, rr.Body.String())
}

func TestCreateMessage_InternalServerError(t *testing.T) {
	services.DeviceService = &deviceMock{}
	createDevice = func(message *models.Device) (*models.DeviceWithId, models.MessageErr) {
		return &models.DeviceWithId{}, models.NewInternalServerError("database error")
	}
	inputJson := `{"name": "name", "interval": 5, "value": 1.4}`
	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandlePost)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"message":"database error","status":500,"error":"server_error"}`, rr.Body.String())
}
