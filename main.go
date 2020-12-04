package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	m := make(chan Measurement)
	mw := MeasurementPrintWriter{measurement: m}
	go mw.ReadMeasurement()

	deviceService := &DeviceService{NewInMemRepo()}
	deviceController := DeviceController{deviceService}
	tickerController := TickerController{&TickerService{DeviceService: *deviceService, measurement: m, stop: make(chan struct{}), Ticker: &MeasurementTicker{}}}
	r := SetupRouter(deviceController, tickerController)

	port := os.Getenv("TSM_PORT")
	if port == "" {
		port = "8000"
	}
	log.Print("TSM | HTTP server listening on :", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
