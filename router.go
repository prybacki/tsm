package main

import (
	"github.com/gorilla/mux"
)

type DeviceController struct {
	DeviceService DeviceCreator
}

func SetupRouter(r *mux.Router) {
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	r.HandleFunc("/devices", deviceController.HandlePost).Methods("POST")
}
