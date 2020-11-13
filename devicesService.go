package main

type DeviceSaver interface {
	Save(*Device) (*DeviceWithId, error)
}
type DeviceService struct {
	DeviceRepo DeviceSaver
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
