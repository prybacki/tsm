package main

import (
	"sync"
)

type inMemoryRepository struct {
	mu      sync.Mutex
	devices map[string]DeviceWithId
	keys    []string
}

func NewInMemRepo() DeviceRepo {
	deviceRepo := inMemoryRepository{devices: make(map[string]DeviceWithId)}
	return &deviceRepo
}

func (r *inMemoryRepository) Save(device *Device, id string) (*DeviceWithId, error) {
	deviceWithId := DeviceWithId{Device: device}
	deviceWithId.Id = id
	r.mu.Lock()
	r.keys = append(r.keys, id)
	r.devices[deviceWithId.Id] = deviceWithId
	r.mu.Unlock()
	return &deviceWithId, nil
}

func (r *inMemoryRepository) GetById(id string) (*DeviceWithId, error) {
	if device, ok := r.devices[id]; ok {
		return &device, nil
	}
	return nil, nil
}

func (r *inMemoryRepository) Get(limit int, page int) ([]DeviceWithId, error) {
	start := limit * page
	end := limit * (page + 1)
	if start > len(r.keys) {
		return []DeviceWithId{}, nil
	}

	var k []string
	switch {
	case limit == 0:
		k = r.keys[:]
	case end > len(r.keys):
		k = r.keys[start:]
	default:
		k = r.keys[start:end]
	}

	v := make([]DeviceWithId, 0)
	for _, value := range k {
		v = append(v, r.devices[value])
	}
	return v, nil
}
