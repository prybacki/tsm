package main

import "fmt"

type MeasurementPrintWriter struct {
	measurement chan Measurement
}

func (mw *MeasurementPrintWriter) ReadMeasurement() error {
	for {
		m := <-mw.measurement
		fmt.Println("Id:", m.Id, "Value: ", m.Value)
	}
}
