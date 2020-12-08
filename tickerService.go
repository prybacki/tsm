package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"sync"
	"time"
)

type Ticker interface {
	Tick(device DeviceWithId, tickerService TickerService)
}

type TickerService struct {
	DeviceService DeviceService
	channel       *amqp.Channel
	isRunning     bool
	stop          chan struct{}
	mu            sync.Mutex
	Ticker
}

type MeasurementTicker struct{}

type Measurement struct {
	Id    int
	Value float64
}

func (ts *TickerService) Start() (started bool, error error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if !ts.isRunning {
		ts.stop = make(chan struct{})
		ts.isRunning = true
		deviceWithId, err := ts.DeviceService.Get(0, 0)
		if err != nil {
			return false, err
		}
		for _, device := range deviceWithId {
			go ts.Ticker.Tick(device, *ts)
		}
		return true, nil
	}
	return false, nil
}

func (ts *TickerService) Stop() (stopped bool, error error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if ts.isRunning {
		ts.isRunning = false
		close(ts.stop)
		return true, nil
	}
	return false, nil
}

func (mt *MeasurementTicker) Tick(device DeviceWithId, ts TickerService) {
	t := time.NewTicker(time.Duration(device.Interval) * time.Second)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			err := ts.channel.Publish(
				"measurement",           // exchange
				strconv.Itoa(device.Id), // routing key
				false,                   // mandatory
				false,                   // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(fmt.Sprintf("%f", device.Value)),
				})
			if err != nil {
				log.Println(err)
			}
		case <-ts.stop:
			return
		}
	}
}
