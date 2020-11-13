package main

import (
	"time"
)

type TickerService struct {
	DeviceService DeviceService
	measurement   chan Measurement
	isRunning     bool
}

type Measurement struct {
	Id    int
	Value float64
}

func (ts *TickerService) Start() (started bool, error error) {
	if !ts.isRunning {
		deviceWithId, err := ts.DeviceService.Get(0, 0)
		if err != nil {
			return false, err
		}
		for _, device := range deviceWithId {
			ticker := time.NewTicker(time.Duration(device.Interval) * time.Second)
			go func() {
				for range ticker.C {
					ts.measurement <- Measurement{Id: device.Id, Value: float64(device.Value)}
				}
			}()
		}
		return true, nil
	}
	return false, nil
}
