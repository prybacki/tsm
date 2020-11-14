package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type deviceService interface {
	Create(*Device) (*DeviceWithId, error)
	Get(int) (*DeviceWithId, error)
}

type DeviceController struct {
	DeviceService deviceService
}

func (dc *DeviceController) HandlePost(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Location", r.Host+r.URL.Path+"/"+strconv.Itoa(deviceWithId.Id))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deviceWithId)
}

func (dc *DeviceController) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	id := mux.Vars(r)["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(NewBadRequestError("id must be a string"))
		return
	}

	deviceWithId, err := dc.DeviceService.Get(idInt)
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
