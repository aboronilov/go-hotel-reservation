package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	SeaSideRoomType
	DeluxeRoomType
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type" json:"type"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Price     float64            `bson:"price" json:"price"`
	HotelID   primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}

type UpdateRoomParams struct {
	BasePrice float64 `bson:"basePrice" json:"basePrice"`
	Price     float64 `bson:"price" json:"price"`
}

func (p *UpdateRoomParams) ToBson() bson.M {
	m := bson.M{}
	if p.BasePrice > 0 {
		m["basePrice"] = p.BasePrice
	}
	if p.Price > 0 {
		m["price"] = p.Price
	}
	return m
}
