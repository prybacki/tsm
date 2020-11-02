package repositories

import (
	"sync/atomic"
	"tsm/models"
)

var (
	DeviceRepo deviceRepoInterface = &inMemoryRepository{devices: make(map[int]models.DeviceWithId)}
)

type deviceRepoInterface interface {
	SaveDevice(*models.Device) (*models.DeviceWithId, error)
}

type inMemoryRepository struct {
	id      int32
	devices map[int]models.DeviceWithId
}

func (r *inMemoryRepository) SaveDevice(device *models.Device) (*models.DeviceWithId, error) {
	deviceWithId := models.DeviceWithId{Device: device}
	deviceWithId.Id = int(atomic.AddInt32(&r.id, 1))
	r.devices[deviceWithId.Id] = deviceWithId
	return &deviceWithId, nil
}
