package main

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	influxCS := "http://" + os.Getenv("INFLUX_HOST") + ":" + os.Getenv("INFLUX_PORT")
	client := influxdb2.NewClientWithOptions(influxCS, "", influxdb2.DefaultOptions().SetPrecision(time.Second))
	if _, err := client.Health(context.Background()); err != nil {
		log.Print(err)
	}
	writeAPI := client.WriteAPIBlocking("", os.Getenv("INFLUX_DB"))

	mongoCS := "mongodb://" + os.Getenv("MONGO_HOST") + ":" + os.Getenv("MONGO_PORT")
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoCS))
	if err != nil {
		log.Print(err)
	}

	m := make(chan Measurement)
	mw := MeasurementWriter{writeAPI, m}
	go mw.Start()

	deviceService := &DeviceService{&MongoDbRepository{*mongoClient}}
	deviceController := DeviceController{deviceService}
	tickerController := TickerController{&TickerService{DeviceService: *deviceService, measurement: m, Ticker: &MeasurementTicker{}}}
	r := SetupRouter(deviceController, tickerController)

	port := os.Getenv("TSM_PORT")
	if port == "" {
		port = "8000"
	}
	log.Print("TSM | HTTP server listening on :", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
