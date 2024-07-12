package db

import (
	"context"

	"github.com/aboronilov/go-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	GetRoomByID(context.Context, string) (*types.Room, error)
	GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error)
	CreateRoom(context.Context, *types.Room) (*types.Room, error)
	DeleteRoomByID(context.Context, string) error
	UpdateRoomByID(ctx context.Context, filter bson.M, params types.UpdateRoomParams) error
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBNAME).Collection(ROOM_COLLECTION),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) CreateRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = res.InsertedID.(primitive.ObjectID)

	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.ID}}
	if err := s.HotelStore.UpdateHotelByID(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}

func (s *MongoRoomStore) GetRoomByID(ctx context.Context, id string) (*types.Room, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var room types.Room
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&room); err != nil {
		return nil, err
	}

	return &room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var rooms []*types.Room
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) DeleteRoomByID(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	if _, err := s.coll.DeleteOne(ctx, filter); err != nil {
		return err
	}

	hotel, _ := s.HotelStore.GetHotelByID(ctx, oid)
	filter = bson.M{"_id": hotel.ID}
	update := bson.M{"$pull": bson.M{"rooms": oid}}
	if err := s.HotelStore.UpdateHotelByID(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (s *MongoRoomStore) UpdateRoomByID(ctx context.Context, filter bson.M, params types.UpdateRoomParams) error {
	update := bson.M{
		"$set": params.ToBson(),
	}
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
