package main

import (
	"context"
	"log"

	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/aboronilov/go-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(firstName, lastName, email string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  "secret",
	})
	if err != nil {
		log.Fatal(err)
	}

	newUser, err := userStore.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created user: %+v\n", newUser)
}

func seedHotel(name, location string, rating int) {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	insertedHotel, err := hotelStore.CreateHotel(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}

	room_1 := &types.Room{
		HotelID: insertedHotel.ID,
		Type:    1,
		Size:    "small",
		Price:   100,
	}
	room_2 := &types.Room{
		HotelID: insertedHotel.ID,
		Type:    2,
		Size:    "normal",
		Price:   120,
	}
	room_3 := &types.Room{
		HotelID: insertedHotel.ID,
		Type:    3,
		Size:    "kingsize",
		Price:   140,
	}
	rooms := []*types.Room{room_1, room_2, room_3}
	for _, room := range rooms {
		insertedRoom, err := roomStore.CreateRoom(ctx, room)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Created room: %+v\n", insertedRoom)
	}

	log.Printf("Created hotel: %+v\n", insertedHotel)
}

func main() {
	seedHotel("IBIS", "New York", 5)
	seedHotel("Tamara Hotel", "Amsterdam", 3)
	seedHotel("Holiday INN", "Paris", 4)
	seedUser("Jack", "Bauer", "jack_bauer@ctu.com")
}

func init() {
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Collection(db.HOTEL_COLLECTION).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Collection(db.ROOM_COLLECTION).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
