package main

import (
	"github.com/stretchr/testify/assert"
	"math"
	"sync"
	"testing"
)

type mockTicker struct {
	wg sync.WaitGroup
}

func (m *mockTicker) Tick(DeviceWithId, TickerService) {
	m.wg.Done()
}

var device = &Device{
	Name:     "test device",
	Interval: 1,
	Value:    2.3,
}

func TestTickThreadIsCreatedForEachDevice(t *testing.T) {
	mockTicker := &mockTicker{}
	deviceService := &DeviceService{NewInMemRepo()}
	sut := &TickerService{DeviceService: *deviceService, Ticker: mockTicker}
	deviceService.Create(device)
	mockTicker.wg.Add(1)
	deviceService.Create(device)
	mockTicker.wg.Add(1)

	sut.Start()

	mockTicker.wg.Wait()
}

func TestDeviceSentAtLeastOneMeasurement(t *testing.T) {
	m := make(chan Measurement)
	deviceService := &DeviceService{NewInMemRepo()}
	sut := &TickerService{DeviceService: *deviceService, measurement: m, Ticker: &MeasurementTicker{}}
	device := &Device{
		Name:     "test device",
		Interval: 1,
		Value:    2.3,
	}
	deviceService.Create(device)

	sut.Start()

	assert.EqualValues(t, 2.3, math.Round((<-m).Value*10)/10)
}
