package main

import (
	"github.com/gorilla/mux"
)

func SetupRouter(deviceController DeviceController, tickerController TickerController) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/devices", deviceController.HandleDevicesPost).Methods("POST")
	r.HandleFunc("/devices/{id}", deviceController.HandleDeviceGet).Methods("GET")
	r.HandleFunc("/devices", deviceController.HandleDevicesGet).Methods("GET")

	r.HandleFunc("/start", tickerController.HandleStart).Methods("POST")
	r.HandleFunc("/stop", tickerController.HandleStop).Methods("POST")
	return r
}
