package main

const (
	badRequest  = "bad_request"
	serverError = "server_error"
)

type Device struct {
	Name     string  `json:"name"`
	Interval int     `json:"interval"`
	Value    float32 `json:"value"`
}

type DeviceWithId struct {
	Id int `json:"id"`
	*Device
}

type MessageErr struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewBadRequestError(message string) *MessageErr {
	return &MessageErr{
		Message: message,
		Error:   badRequest,
	}
}

func NewInternalServerError(message string) *MessageErr {
	return &MessageErr{
		Message: message,
		Error:   serverError,
	}
}

func (d *Device) Validate() *MessageErr {
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
