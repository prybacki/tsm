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

	r := SetupRouter(m)
	port := os.Getenv("TSM_PORT")
	if port == "" {
		port = "8000"
	}
	log.Print("TSM | HTTP server listening on :", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
