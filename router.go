package main

import (
	"github.com/gorilla/mux"
)

func SetupRouter(m chan Measurement) *mux.Router {
	r := mux.NewRouter()
	deviceService := &DeviceService{NewInMemRepo()}
	deviceController := DeviceController{deviceService}
	r.HandleFunc("/devices", deviceController.HandleDevicesPost).Methods("POST")
	r.HandleFunc("/devices/{id}", deviceController.HandleDeviceGet).Methods("GET")
	r.HandleFunc("/devices", deviceController.HandleDevicesGet).Methods("GET")

	tickerController := TickerController{&TickerService{DeviceService: *deviceService, measurement: m}}
	r.HandleFunc("/start", tickerController.HandleTickerPost).Methods("POST")
	return r
}
