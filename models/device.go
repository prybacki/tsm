package models

type Device struct {
	Name     string  `json:"name"`
	Interval int     `json:"interval"`
	Value    float32 `json:"value"`
}

type DeviceWithId struct {
	Id int `json:"id"`
	*Device
}

func (d *Device) Validate() MessageErr {
	var msg string
	if d.Interval <= 0 {
		msg += "interval has to be greater than 0;"
	}
	if d.Name == "" {
		msg += "device name can't be empty;"
	}

	if msg != "" {
		return NewBadRequestError("Request cannot be validate: " + msg)
	}
	return nil
}
