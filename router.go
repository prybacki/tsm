package main

import (
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	r.HandleFunc("/devices", deviceController.HandleDevicesPost).Methods("POST")
	r.HandleFunc("/devices/{id}", deviceController.HandleDeviceGet).Methods("GET")
	r.HandleFunc("/devices", deviceController.HandleDevicesGet).Methods("GET")
	return r
}
