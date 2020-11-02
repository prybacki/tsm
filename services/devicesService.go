package services

import (
	"tsm/models"
	"tsm/repositories"
)

var (
	DeviceService deviceServiceInterface = &deviceService{}
)

type deviceServiceInterface interface {
	CreateDevice(*models.Device) (*models.DeviceWithId, models.MessageErr)
}

type deviceService struct{}

func (ds *deviceService) CreateDevice(device *models.Device) (*models.DeviceWithId, models.MessageErr) {
	if err := device.Validate(); err != nil {
		return &models.DeviceWithId{}, err
	}
	d, err := repositories.DeviceRepo.SaveDevice(device)
	if err != nil {
		return &models.DeviceWithId{}, models.NewInternalServerError("database error")
	}
	return d, nil
}
