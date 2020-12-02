package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	connectionString := "http://" + os.Getenv("INFLUX_HOST") + ":" + os.Getenv("INFLUX_PORT")
	client := influxdb2.NewClientWithOptions(connectionString, "", influxdb2.DefaultOptions().SetPrecision(time.Second))
	if _, err := client.Health(context.Background()); err != nil {
		panic("Cannot connect to influxdb")
	}
	writeAPI := client.WriteAPIBlocking("", os.Getenv("INFLUX_DB"))

	m := make(chan Measurement)
	mw := MeasurementWriter{writeAPI, m}
	go mw.StoreMeasurement()

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
