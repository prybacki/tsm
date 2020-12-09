package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

const (
	defaultLimit = 100
	defaultPage  = 0
)

type deviceService interface {
	Create(*Device) (*DeviceWithId, error)
	GetById(string) (*DeviceWithId, error)
	Get(int, int) ([]DeviceWithId, error)
}

type DeviceController struct {
	DeviceService deviceService
}

func (dc *DeviceController) HandleDevicesPost(w http.ResponseWriter, r *http.Request) {
	var device Device
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		e := NewBadRequestError("invalid json body")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(e)
		return
	}

	deviceWithId, err := dc.DeviceService.Create(&device)
	if err != nil {
		switch err.(*MessageErr).Code {
		case badRequest:
			w.WriteHeader(http.StatusBadRequest)
		case serverError:
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Location", r.Host+r.URL.Path+"/"+deviceWithId.Id)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deviceWithId)
}

func (dc *DeviceController) HandleDeviceGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewBadRequestError("id must be a hex string"))
		return
	}

	deviceWithId, err := dc.DeviceService.GetById(id)
	if err != nil {
		switch err.(*MessageErr).Code {
		case notFound:
			w.WriteHeader(http.StatusNotFound)
		case serverError:
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceWithId)
}

func (dc *DeviceController) HandleDevicesGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	limit, page := defaultLimit, defaultPage
	lErr := dc.readIntQueryParam(r, "limit", &limit)
	pErr := dc.readIntQueryParam(r, "page", &page)
	if lErr != nil || pErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewBadRequestError("given string value in limit or page query parameters"))
		return
	}

	deviceWithId, err := dc.DeviceService.Get(limit, page)
	if err != nil {
		switch err.(*MessageErr).Code {
		case badRequest:
			w.WriteHeader(http.StatusBadRequest)
		case serverError:
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(err)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(deviceWithId)
}

func (dc *DeviceController) readIntQueryParam(r *http.Request, param string, result *int) error {
	paramS := r.URL.Query().Get(param)
	if paramS != "" {
		param, lErr := strconv.Atoi(paramS)
		if lErr != nil {
			return lErr
		}
		*result = param
	}
	return nil
}
