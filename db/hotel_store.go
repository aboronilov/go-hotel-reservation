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
	GetHotels(context.Context, bson.M) ([]*types.Hotel, error)
	CreateHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	DeleteHotelByID(context.Context, string) error
	UpdateHotelByID(ctx context.Context, filter bson.M, params bson.M) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(HOTEL_COLLECTION),
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

func (s *MongoHotelStore) UpdateHotelByID(ctx context.Context, filter bson.M, update bson.M) error {
	// update := bson.M{
	// 	"$set": params.ToBson(),
	// }
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var hotel types.Hotel
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel)
	if err != nil {
		return nil, err
	}

	return &hotel, nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M) ([]*types.Hotel, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	for cur.Next(ctx) {
		var hotel types.Hotel
		err := cur.Decode(&hotel)
		if err != nil {
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}

	return hotels, nil
}

func (s *MongoHotelStore) DeleteHotelByID(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	return nil
}
