package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

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
		switch err.Error {
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
