package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Size    string             `bson:"size" json:"size"`
	Seaside bool               `bson:"seaside" json:"seaside"`
	Price   float64            `bson:"price" json:"price"`
	HotelID primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}

type UpdateRoomParams struct {
	BasePrice float64 `bson:"basePrice" json:"basePrice"`
	Price     float64 `bson:"price" json:"price"`
}

func (p *UpdateRoomParams) ToBson() bson.M {
	m := bson.M{}
	if p.Price > 0 {
		m["price"] = p.Price
	}
	return m
}
