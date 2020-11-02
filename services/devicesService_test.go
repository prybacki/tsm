package services

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"tsm/models"
	"tsm/repositories"
)

var (
	saveDevice func(msg *models.Device) (*models.DeviceWithId, error)
)

type repoMock struct{}

func (m *repoMock) SaveDevice(msg *models.Device) (*models.DeviceWithId, error) {
	return saveDevice(msg)
}

func TestDevicesService_CreateDevice_Success(t *testing.T) {
	device := &models.Device{
		Name:     "test device",
		Interval: 50,
		Value:    2.3,
	}

	d, err := DeviceService.CreateDevice(device)

	assert.Nil(t, err)
	assert.EqualValues(t, 1, d.Id)
	assert.EqualValues(t, "test device", d.Name)
	assert.EqualValues(t, 50, d.Interval)
	assert.EqualValues(t, 2.3, d.Value)
}

func TestDevicesService_CreateDevice_Fail_Interval(t *testing.T) {
	device := &models.Device{
		Name:     "test device",
		Interval: -1,
		Value:    2.3,
	}

	_, err := DeviceService.CreateDevice(device)

	assert.NotNil(t, err)
	assert.EqualValues(t, "Request cannot be validate: interval has to be greater than 0;", err.Message())
	assert.EqualValues(t, 400, err.Status())
	assert.EqualValues(t, "bad_request", err.Error())
}

func TestDevicesService_CreateDevice_Fail_Name(t *testing.T) {
	device := &models.Device{
		Name:     "",
		Interval: 5,
		Value:    2.3,
	}

	_, err := DeviceService.CreateDevice(device)

	assert.NotNil(t, err)
	assert.EqualValues(t, "Request cannot be validate: device name can't be empty;", err.Message())
	assert.EqualValues(t, 400, err.Status())
	assert.EqualValues(t, "bad_request", err.Error())
}

func TestDevicesService_CreateDevice_Database_Error(t *testing.T) {
	repositories.DeviceRepo = &repoMock{}
	saveDevice = func(messageId *models.Device) (*models.DeviceWithId, error) {
		return &models.DeviceWithId{}, errors.New("some error")
	}
	device := &models.Device{
		Name:     "test device",
		Interval: 50,
		Value:    2.3,
	}

	_, err := DeviceService.CreateDevice(device)

	assert.NotNil(t, err)
	assert.EqualValues(t, "database error", err.Message())
	assert.EqualValues(t, 500, err.Status())
	assert.EqualValues(t, "server_error", err.Error())
}
