package main

import (
	"context"
	"log"

	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/aboronilov/go-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client, db.DBNAME)
	roomStore := db.NewMongoRoomStore(client, db.DBNAME)

	hotel := &types.Hotel{
		Name:     "IBIS",
		Location: "New York",
	}
	insertedHotel, err := hotelStore.CreateHotel(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}

	room := &types.Room{
		HotelID:   insertedHotel.ID,
		Type:      1,
		BasePrice: 99.9,
	}
	insertedRoom, err := roomStore.CreateRoom(ctx, room)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created hotel: %+v\n", insertedHotel)
	log.Printf("Created room: %+v\n", insertedRoom)
}
