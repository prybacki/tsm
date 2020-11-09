package main

type DeviceCreator interface {
	Create(*Device) (*DeviceWithId, *MessageErr)
}

type DeviceService struct {
	DeviceRepo DeviceSaver
}

type DeviceSaver interface {
	Save(*Device) (*DeviceWithId, error)
}

func (ds *DeviceService) Create(device *Device) (*DeviceWithId, *MessageErr) {
	if err := device.Validate(); err != nil {
		return &DeviceWithId{}, err
	}
	d, err := ds.DeviceRepo.Save(device)
	if err != nil {
		return &DeviceWithId{}, NewInternalServerError("database error")
	}
	return d, nil
}
