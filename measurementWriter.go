package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"time"
)

type MeasurementWriter struct {
	writeApi api.WriteAPIBlocking
	ch       *amqp.Channel
}

func (mw *MeasurementWriter) Start() {
	msgs, err := mw.ch.Consume(
		"measurement", // queue
		"",            // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		log.Println(err)
		return
	}
	for d := range msgs {
		v, err := strconv.ParseFloat(string(d.Body), 32)
		if err != nil {
			log.Println(err)
		}
		p := influxdb2.NewPointWithMeasurement("deviceValues").
			AddTag("id", d.RoutingKey).
			AddField("value", v).
			SetTime(time.Now())
		err = mw.writeApi.WritePoint(context.Background(), p)
		if err != nil {
			log.Println(err)
		}
	}
}
