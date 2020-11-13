package main

import (
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()
	deviceController := DeviceController{&DeviceService{NewInMemRepo()}}
	r.HandleFunc("/devices", deviceController.HandlePost).Methods("POST")
	return r
}
