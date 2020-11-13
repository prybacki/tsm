package main

import (
	"sync"
)

type inMemoryRepository struct {
	mu      sync.Mutex
	id      int
	devices map[int]DeviceWithId
}

func NewInMemRepo() DeviceSaver {
	deviceRepo := inMemoryRepository{devices: make(map[int]DeviceWithId)}
	return &deviceRepo
}

func (r *inMemoryRepository) Save(device *Device) (*DeviceWithId, error) {
	deviceWithId := DeviceWithId{Device: device}
	r.mu.Lock()
	r.id++
	deviceWithId.Id = r.id
	r.devices[deviceWithId.Id] = deviceWithId
	r.mu.Unlock()
	return &deviceWithId, nil
}