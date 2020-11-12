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
		return &DeviceWithId{}, err
	}
	d, err := ds.DeviceRepo.Save(device)
	if err != nil {
		return &DeviceWithId{}, NewInternalServerError("database error")
	}
	return d, nil
}

func (ds *DeviceService) Get(id int) (*DeviceWithId, error) {
	d, err := ds.DeviceRepo.GetById(id)
	if err != nil {
		return &DeviceWithId{}, NewInternalServerError("database error")
	}
	if d.Device == nil {
		return &DeviceWithId{}, NewNotFoundError("device not found")
	}
	return d, nil
}
