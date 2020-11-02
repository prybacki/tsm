package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tsm/models"
	"tsm/services"
)

func HandlePost(w http.ResponseWriter, r *http.Request) {
	var device models.Device
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		e := models.NewBadRequestError("invalid json body")
		w.WriteHeader(e.Status())
		json.NewEncoder(w).Encode(e)
		return
	}

	deviceWithId, err := services.DeviceService.CreateDevice(&device)
	if err != nil {
		w.WriteHeader(err.Status())
		json.NewEncoder(w).Encode(err)
		return
	}

	w.Header().Set("Location", r.Host+r.URL.Path+"/"+strconv.Itoa(deviceWithId.Id))
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(deviceWithId)
}
