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

	tickerController := TickerController{&TickerService{DeviceService: *deviceService, measurement: m, stop: make(chan struct{}), Ticker: &MeasurementTicker{}}}
	r.HandleFunc("/start", tickerController.HandleStart).Methods("POST")
	r.HandleFunc("/stop", tickerController.HandleStop).Methods("POST")
	return r
}
