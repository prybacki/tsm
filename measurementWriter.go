package main

import (
	"context"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"strconv"
	"time"
)

type MeasurementWriter struct {
	writeApi    api.WriteAPIBlocking
	measurement chan Measurement
}

func (mw *MeasurementWriter) StoreMeasurement() {
	for {
		m := <-mw.measurement
		p := influxdb2.NewPointWithMeasurement("deviceValues").
			AddTag("id", strconv.Itoa(m.Id)).
			AddField("value", m.Value).
			SetTime(time.Now())
		mw.writeApi.WritePoint(context.Background(), p)
	}
}
