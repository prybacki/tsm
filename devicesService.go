package main

import "log"

type DeviceRepo interface {
	Save(*Device) (*DeviceWithId, error)
	GetById(int) (*DeviceWithId, error)
	Get(int, int) ([]DeviceWithId, error)
}

type DeviceService struct {
	DeviceRepo DeviceRepo
}

func (ds *DeviceService) Create(device *Device) (*DeviceWithId, error) {
	if err := device.Validate(); err != nil {
		return nil, err
	}
	d, err := ds.DeviceRepo.Save(device)
	if err != nil {
		log.Println("Error during create device: ", err.Error())
		return nil, NewInternalServerError("database error")
	}
	return d, nil
}

func (ds *DeviceService) GetById(id int) (*DeviceWithId, error) {
	d, err := ds.DeviceRepo.GetById(id)
	if err != nil {
		log.Println("Error during get device by id: ", err.Error())
		return nil, NewInternalServerError("database error")
	}
	if d == nil {
		return nil, NewNotFoundError("device not found")
	}
	return d, nil
}

func (ds *DeviceService) Get(limit int, page int) ([]DeviceWithId, error) {
	if limit < 0 || page < 0 {
		return nil, NewBadRequestError("negative limit or page")
	}
	d, err := ds.DeviceRepo.Get(limit, page)
	if err != nil {
		log.Println("Error during get devices: ", err.Error())
		return nil, NewInternalServerError("database error")
	}
	return d, nil
}
