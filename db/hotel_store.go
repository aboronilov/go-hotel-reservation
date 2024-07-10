package db

import (
	"context"

	"github.com/aboronilov/go-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	GetHotels(context.Context) ([]*types.Hotel, error)
	CreateHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	DeleteHotelByID(context.Context, string) error
	UpdateHotelByID(ctx context.Context, filter bson.M, params types.UpdateHotelParams) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbname string) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection(HOTEL_COLLECTION),
	}
}

func (s *MongoHotelStore) CreateHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	res, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = res.InsertedID.(primitive.ObjectID)

	return hotel, nil
}

func (s *MongoHotelStore) UpdateHotelByID(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {
	update := bson.M{
		"$set": params.ToBson(),
	}
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
