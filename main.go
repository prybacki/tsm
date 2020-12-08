package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	influxCS := "http://" + os.Getenv("INFLUX_HOST") + ":" + os.Getenv("INFLUX_PORT")
	client := influxdb2.NewClientWithOptions(influxCS, "", influxdb2.DefaultOptions().SetPrecision(time.Second))
	if _, err := client.Health(context.Background()); err != nil {
		log.Println(err)
	}
	writeAPI := client.WriteAPIBlocking("", os.Getenv("INFLUX_DB"))

	rabbitCS := "amqp://guest:guest@" + os.Getenv("RABBIT_HOST") + ":" + os.Getenv("RABBIT_PORT")
	conn, err := amqp.Dial(rabbitCS)
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
	}
	defer ch.Close()
	prepareRabbit(ch)

	mw := MeasurementWriter{writeAPI, ch}
	go mw.Start()

	deviceService := &DeviceService{NewInMemRepo()}
	deviceController := DeviceController{deviceService}
	tickerController := TickerController{&TickerService{DeviceService: *deviceService, channel: ch, Ticker: &MeasurementTicker{}}}
	r := SetupRouter(deviceController, tickerController)

	port := os.Getenv("TSM_PORT")
	if port == "" {
		port = "8000"
	}
	log.Println("TSM | HTTP server listening on :", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func prepareRabbit(ch *amqp.Channel) {
	err := ch.ExchangeDeclare(
		"measurement", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = ch.QueueDeclare(
		"measurement", // name
		false,         // durable
		false,         // delete when unused
		true,          // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Println(err)
		return
	}

	err = ch.QueueBind(
		"measurement", // queue name
		"#",           // routing key
		"measurement", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Println(err)
		return
	}
}
