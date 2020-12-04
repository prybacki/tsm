package main

const (
	badRequest  = "bad_request"
	serverError = "server_error"
	notFound    = "not_found"
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

type Message struct {
	Message string `json:"message"`
}

type MessageErr struct {
	Message string `json:"message"`
	Code    string `json:"error"`
}

func (e *MessageErr) Error() string { return e.Message }

func NewBadRequestError(message string) *MessageErr {
	return &MessageErr{
		Message: message,
		Code:    badRequest,
	}
}

func NewInternalServerError(message string) *MessageErr {
	return &MessageErr{
		Message: message,
		Code:    serverError,
	}
}

func NewNotFoundError(message string) *MessageErr {
	return &MessageErr{
		Message: message,
		Code:    notFound,
	}
}

func NewMessage(message string) *Message {
	return &Message{
		Message: message,
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
