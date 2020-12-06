package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

type repoMock struct {
	returnValue []DeviceWithId
	error
}

func (r repoMock) Save(*Device, string) (*DeviceWithId, error) {
	if r.returnValue == nil {
		return nil, r.error
	}
	return &(r.returnValue)[0], r.error
}

func (r repoMock) GetById(string) (*DeviceWithId, error) {
	if r.returnValue == nil {
		return nil, r.error
	}
	return &(r.returnValue)[0], r.error
}

func (r repoMock) Get(int, int) ([]DeviceWithId, error) {
	return r.returnValue, r.error
}

//test Create
func TestCreateDevice_Success(t *testing.T) {
	sut := DeviceService{NewInMemRepo()}
	device := &Device{
		Name:     "test device",
		Interval: 50,
		Value:    2.3,
	}

	d, err := sut.Create(device)

	assert.Nil(t, err)
	resultId, err := primitive.ObjectIDFromHex(d.Id)
	assert.NoError(t, err)
	assert.IsType(t, primitive.ObjectID{}, resultId)
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

//test GetById
func TestGetDevice_Success(t *testing.T) {
	id := primitive.NewObjectID().Hex()
	repo := NewInMemRepo()
	repo.Save(&Device{Name: "name", Interval: 5, Value: 1.4}, id)
	sut := DeviceService{repo}

	d, err := sut.GetById(id)

	assert.Nil(t, err)
	assert.EqualValues(t, id, d.Id)
	assert.EqualValues(t, "name", d.Name)
	assert.EqualValues(t, 5, d.Interval)
	assert.EqualValues(t, 1.4, d.Value)
}

func TestGetDevice_DatabaseError(t *testing.T) {
	repoMock := repoMock{error: errors.New("some error")}
	sut := DeviceService{repoMock}

	_, err := sut.GetById(primitive.NewObjectID().Hex())

	assert.NotNil(t, err)
	assert.EqualValues(t, "database error", err.Error())
	assert.EqualValues(t, "server_error", err.(*MessageErr).Code)
}

func TestGetDevice_NotFoundError(t *testing.T) {
	repoMock := repoMock{returnValue: nil}
	sut := DeviceService{repoMock}

	_, err := sut.GetById(primitive.NewObjectID().Hex())

	assert.NotNil(t, err)
	assert.EqualValues(t, "device not found", err.Error())
	assert.EqualValues(t, "not_found", err.(*MessageErr).Code)
}

//test Get

var limitAndPageTests = []struct {
	limit         int
	page          int
	expectedSize  int
	expectedNames []string
}{
	{0, 0, 3, []string{"name1", "name2", "name3"}},
	{0, 1, 3, []string{"name1", "name2", "name3"}},
	{0, 2, 3, []string{"name1", "name2", "name3"}},
	{0, 3, 3, []string{"name1", "name2", "name3"}},

	{1, 0, 1, []string{"name1"}},
	{1, 1, 1, []string{"name2"}},
	{1, 2, 1, []string{"name3"}},
	{1, 3, 0, nil},

	{2, 0, 2, []string{"name1", "name2"}},
	{2, 1, 1, []string{"name3"}},
	{2, 2, 0, nil},

	{3, 0, 3, []string{"name1", "name2", "name3"}},
	{3, 1, 0, nil},

	{4, 0, 3, []string{"name1", "name2", "name3"}},
	{4, 1, 0, nil},
}

func TestGetDevices_Success(t *testing.T) {
	repo := NewInMemRepo()
	repo.Save(&Device{Name: "name1"}, primitive.NewObjectID().Hex())
	repo.Save(&Device{Name: "name2"}, primitive.NewObjectID().Hex())
	repo.Save(&Device{Name: "name3"}, primitive.NewObjectID().Hex())
	sut := DeviceService{repo}

	for _, v := range limitAndPageTests {
		d, err := sut.Get(v.limit, v.page)
		assert.Nil(t, err)
		assert.EqualValues(t, v.expectedSize, len(d))
		assert.EqualValues(t, v.expectedNames, nameList(d))
	}
}

func nameList(devices []DeviceWithId) []string {
	var result []string
	for _, d := range devices {
		result = append(result, d.Name)
	}
	return result
}
