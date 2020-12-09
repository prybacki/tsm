package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDbRepository struct {
	client mongo.Client
}

func (r *MongoDbRepository) Save(device *Device) (*DeviceWithId, error) {
	deviceWithId := DeviceWithId{Device: device}
	deviceWithId.Id = primitive.NewObjectID().Hex()

	col := r.client.Database("tsm").Collection("devices")
	_, err := col.InsertOne(context.Background(), deviceWithId)
	if err != nil {
		return nil, err
	}

	return &deviceWithId, nil
}

func (r *MongoDbRepository) GetById(id string) (*DeviceWithId, error) {
	var device DeviceWithId
	col := r.client.Database("tsm").Collection("devices")
	err := col.FindOne(context.Background(), bson.M{"_id": id}).Decode(&device)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, nil
		}
		return nil, err
	}
	return &device, nil
}

func (r *MongoDbRepository) Get(limit int, page int) ([]DeviceWithId, error) {
	findOptions := options.Find()
	if limit != 0 {
		findOptions.SetLimit(int64(limit))
		findOptions.SetSkip(int64(page * limit))
	}

	col := r.client.Database("tsm").Collection("devices")
	cur, err := col.Find(context.Background(), bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}

	v := make([]DeviceWithId, 0)
	for cur.Next(context.Background()) {
		var elem DeviceWithId
		if err := cur.Decode(&elem); err != nil {
			return nil, err
		}

		v = append(v, elem)
	}
	cur.Close(context.Background())
	return v, nil
}
