package main

import (
	"sync"
	"time"
)

type Ticker interface {
	Tick(device DeviceWithId, tickerService TickerService)
}

type TickerService struct {
	DeviceService DeviceService
	measurement   chan Measurement
	isRunning     bool
	stop          chan struct{}
	mu            sync.Mutex
	Ticker
}

type MeasurementTicker struct{}

type Measurement struct {
	Id    string
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
			ts.measurement <- Measurement{Id: device.Id, Value: float64(device.Value)}
		case <-ts.stop:
			return
		}
	}
}
