package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

type repoMock struct {
	returnValue *DeviceWithId
	error
}

func (r repoMock) Save(message *Device) (*DeviceWithId, error) {
	return r.returnValue, r.error
}

func TestDevicesService_CreateDevice_Success(t *testing.T) {
	sut := DeviceService{NewInMemRepo()}
	device := &Device{
		Name:     "test device",
		Interval: 50,
		Value:    2.3,
	}

	d, err := sut.Create(device)

	assert.Nil(t, err)
	assert.EqualValues(t, 1, d.Id)
	assert.EqualValues(t, "test device", d.Name)
	assert.EqualValues(t, 50, d.Interval)
	assert.EqualValues(t, 2.3, d.Value)
}

func TestDevicesService_CreateDevice_Fail_Interval(t *testing.T) {
	sut := DeviceService{NewInMemRepo()}
	device := &Device{
		Name:     "test device",
		Interval: -1,
		Value:    2.3,
	}

	_, err := sut.Create(device)

	assert.NotNil(t, err)
	assert.EqualValues(t, "Request cannot be validate: interval has to be greater than 0;", err.Error())
	assert.EqualValues(t, "bad_request", err.(*MessageErr).Code)
}

func TestDevicesService_CreateDevice_Fail_Name(t *testing.T) {
	sut := DeviceService{NewInMemRepo()}
	device := &Device{
		Name:     "",
		Interval: 5,
		Value:    2.3,
	}

	_, err := sut.Create(device)

	assert.NotNil(t, err)
	assert.EqualValues(t, "Request cannot be validate: device name can't be empty;", err.Error())
	assert.EqualValues(t, "bad_request", err.(*MessageErr).Code)
}

func TestDevicesService_CreateDevice_Database_Error(t *testing.T) {
	repoMock := repoMock{error: errors.New("some error")}
	sut := DeviceService{repoMock}

	device := &Device{
		Name:     "test device",
		Interval: 50,
		Value:    2.3,
	}

	_, err := sut.Create(device)

	assert.NotNil(t, err)
	assert.EqualValues(t, "database error", err.Error())
	assert.EqualValues(t, "server_error", err.(*MessageErr).Code)
}
