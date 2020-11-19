package main

type DeviceRepo interface {
	Save(*Device) (*DeviceWithId, error)
	GetById(int) (*DeviceWithId, error)
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
		return nil, NewInternalServerError("database error")
	}
	return d, nil
}

func (ds *DeviceService) Get(id int) (*DeviceWithId, error) {
	d, err := ds.DeviceRepo.GetById(id)
	if err != nil {
		return nil, NewInternalServerError("database error")
	}
	if d == nil {
		return nil, NewNotFoundError("device not found")
	}
	return d, nil
}
