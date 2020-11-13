package main

import "fmt"

type measurementWriter interface {
	ReadMeasurement() error
}

type MeasurementPrintWriter struct {
	mw          measurementWriter
	measurement chan Measurement
}

func (mw *MeasurementPrintWriter) ReadMeasurement() error {
	for {
		fmt.Println("Id:", (<-mw.measurement).Id, "Value: ", (<-mw.measurement).Value)
	}
}
