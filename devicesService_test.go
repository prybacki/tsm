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

func (r repoMock) GetById(id int) (*DeviceWithId, error) {
	return r.returnValue, r.error
}

func TestCreateDevice_Success(t *testing.T) {
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

func TestCreateDevice_FailInterval(t *testing.T) {
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

func TestCreateDevice_FailName(t *testing.T) {
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

func TestCreateDevice_DatabaseError(t *testing.T) {
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

//
func TestGetDevice_Success(t *testing.T) {
	repo := NewInMemRepo()
	repo.Save(&Device{Name: "name", Interval: 5, Value: 1.4})
	sut := DeviceService{repo}

	d, err := sut.Get(1)

	assert.Nil(t, err)
	assert.EqualValues(t, 1, d.Id)
	assert.EqualValues(t, "name", d.Name)
	assert.EqualValues(t, 5, d.Interval)
	assert.EqualValues(t, 1.4, d.Value)
}

func TestGetDevice_DatabaseError(t *testing.T) {
	repoMock := repoMock{error: errors.New("some error")}
	sut := DeviceService{repoMock}

	_, err := sut.Get(1)

	assert.NotNil(t, err)
	assert.EqualValues(t, "database error", err.Error())
	assert.EqualValues(t, "server_error", err.(*MessageErr).Code)
}

func TestGetDevice_NotFoundError(t *testing.T) {
	repoMock := repoMock{returnValue: nil}
	sut := DeviceService{repoMock}

	_, err := sut.Get(1)

	assert.NotNil(t, err)
	assert.EqualValues(t, "device not found", err.Error())
	assert.EqualValues(t, "not_found", err.(*MessageErr).Code)
}
