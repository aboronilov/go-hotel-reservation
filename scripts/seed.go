package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/aboronilov/go-hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	collections := []string{db.HOTEL_COLLECTION, db.ROOM_COLLECTION, db.USERS_COLLECTION}
	for _, collection := range collections {
		err := client.Database(db.DBNAME).Collection(collection).Drop(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}

	hotelStore := db.NewMongoHotelStore(client)

	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Hotel:   db.NewMongoHotelStore(client),
		Booking: db.NewMongoBookingStore(client),
	}

	newAdmin := fixtures.AddUser(store, "Jack", "Bauer", true)
	fmt.Println("admin --->", newAdmin.ID)

	newUser := fixtures.AddUser(store, "Tony", "Almeida", false)
	fmt.Println("user --->", newUser.ID)

	newHotel := fixtures.AddHotel(store, "IBIS", "New York", 5)
	fmt.Println("hotel --->", newHotel.ID)

	from := time.Now().AddDate(0, 0, 1)
	till := time.Now().AddDate(0, 0, 6)
	newBooking := fixtures.AddBooking(store, newAdmin.ID, newHotel.Rooms[0], from, till)
	fmt.Println("booking --->", newBooking.ID)
}
