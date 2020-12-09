package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
	"testing"
)

type serviceMock struct {
	returnValue []DeviceWithId
	returnError error
}

func (s *serviceMock) Create(*Device) (*DeviceWithId, error) {
	if s.returnValue == nil {
		return nil, s.returnError
	}
	return &(s.returnValue)[0], s.returnError
}

func (s *serviceMock) GetById(string) (*DeviceWithId, error) {
	if s.returnValue == nil {
		return nil, s.returnError
	}
	return &(s.returnValue)[0], s.returnError

}

func (s *serviceMock) Get(int, int) ([]DeviceWithId, error) {
	return s.returnValue, s.returnError

}

//test POST /devices
func TestCreateDevice_Pass(t *testing.T) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	inputJson := `{"name": "name", "interval": 5, "value": 1.4}`
	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleDevicesPost)
	handler.ServeHTTP(rr, req)

	d := DeviceWithId{}
	json.Unmarshal([]byte(rr.Body.String()), &d)
	resultId, err := primitive.ObjectIDFromHex(d.Id)
	assert.EqualValues(t, http.StatusCreated, rr.Code)
	assert.NoError(t, err)
	assert.IsType(t, primitive.ObjectID{}, resultId)
	assert.EqualValues(t, "name", d.Name)
	assert.EqualValues(t, 5, d.Interval)
	assert.EqualValues(t, 1.4, d.Value)

}

func TestCreateDevice_InvalidJson(t *testing.T) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	inputJson := `{"name": 1}`
	req, err := http.NewRequest(http.MethodPost, "/devices", bytes.NewBufferString(inputJson))
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleDevicesPost)
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
	handler := http.HandlerFunc(deviceController.HandleDevicesPost)
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
	handler := http.HandlerFunc(deviceController.HandleDevicesPost)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"message":"database error","error":"server_error"}`, rr.Body.String())
}

//test GET /devices/{id}
func TestGetDevice_Pass(t *testing.T) {
	repo := NewInMemRepo()
	d, _ := repo.Save(&Device{Name: "name", Interval: 5, Value: 1.4})
	deviceController := DeviceController{&DeviceService{repo}}
	req, err := http.NewRequest(http.MethodGet, "/devices", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": d.Id,
	})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleDeviceGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusOK, rr.Code)
	require.JSONEq(t, "{\"id\":\""+d.Id+"\", \"interval\":5, \"name\":\"name\", \"value\":1.4}", rr.Body.String())
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
	handler := http.HandlerFunc(deviceController.HandleDeviceGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"bad_request", "message":"id must be a hex string"}`, rr.Body.String())
}

func TestGetDevice_NotFound(t *testing.T) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	req, err := http.NewRequest(http.MethodGet, "/devices", nil)
	req = mux.SetURLVars(req, map[string]string{
		"id": primitive.NewObjectID().Hex(),
	})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleDeviceGet)
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
		"id": primitive.NewObjectID().Hex(),
	})
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleDeviceGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"message":"database error","error":"server_error"}`, rr.Body.String())
}

//test GET /devices
func TestGetDevices_Pass(t *testing.T) {
	repo := NewInMemRepo()
	d1, _ := repo.Save(&Device{Name: "name1", Interval: 5, Value: 1.4})
	d2, _ := repo.Save(&Device{Name: "name2", Interval: 10, Value: 2.4})
	deviceController := DeviceController{&DeviceService{repo}}
	req, err := http.NewRequest(http.MethodGet, "/devices", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleDevicesGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusOK, rr.Code)
	require.JSONEq(t, "[{\"id\":\""+d1.Id+"\", \"interval\":5, \"name\":\"name1\", \"value\":1.4}, {\"id\":\""+d2.Id+"\", \"interval\":10, \"name\":\"name2\", \"value\":2.4}]", rr.Body.String())
}

func TestGetDevices_Pass_Empty(t *testing.T) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	req, err := http.NewRequest(http.MethodGet, "/devices", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleDevicesGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusOK, rr.Code)
	require.JSONEq(t, `[]`, rr.Body.String())
}

func TestGetDevices_BadRequest_Negative(t *testing.T) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	req, err := http.NewRequest(http.MethodGet, "/devices", nil)
	q := req.URL.Query()
	q.Add("limit", "-1")
	req.URL.RawQuery = q.Encode()
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleDevicesGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"bad_request", "message":"negative limit or page"}`, rr.Body.String())
}

func TestGetDevices_BadRequest_String(t *testing.T) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	req, err := http.NewRequest(http.MethodGet, "/devices", nil)
	q := req.URL.Query()
	q.Add("limit", "string")
	req.URL.RawQuery = q.Encode()
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleDevicesGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
	require.JSONEq(t, `{"error":"bad_request", "message":"given string value in limit or page query parameters"}`, rr.Body.String())
}

func TestGetDevices_InternalServerError(t *testing.T) {
	deviceService := &serviceMock{returnError: NewInternalServerError("database error")}
	deviceController := DeviceController{deviceService}

	req, err := http.NewRequest(http.MethodGet, "/devices", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deviceController.HandleDevicesGet)
	handler.ServeHTTP(rr, req)

	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)
	require.JSONEq(t, `{"message":"database error","error":"server_error"}`, rr.Body.String())
}
