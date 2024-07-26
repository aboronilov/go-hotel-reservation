package db

import (
	"context"

	"github.com/aboronilov/go-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	BookRoom(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingByID(context.Context, primitive.ObjectID) (*types.Booking, error)
	UpdateBooking(context.Context, primitive.ObjectID, bson.M) error
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(DBNAME).Collection(BOOKING_COLLECTION),
	}
}

func (s *MongoBookingStore) BookRoom(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	res, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = res.InsertedID.(primitive.ObjectID)

	return booking, nil
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var bookings []*types.Booking
	if err := cur.All(ctx, &bookings); err != nil {
		return nil, err
	}

	return bookings, nil
}

func (s *MongoBookingStore) GetBookingByID(ctx context.Context, oid primitive.ObjectID) (*types.Booking, error) {
	var booking types.Booking
	err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking)
	if err != nil {
		return nil, err
	}

	return &booking, nil
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}
