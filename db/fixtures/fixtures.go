package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/aboronilov/go-hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddUser(store *db.Store, firstName, lastName string, isAdmin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     fmt.Sprintf("%s_%s@ctu.com", firstName, lastName),
		Password:  fmt.Sprintf("%s_%s", firstName, lastName),
	})
	if err != nil {
		log.Fatal(err)
	}

	user.IsAdmin = isAdmin
	newUser, err := store.User.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}

	return newUser
}

func AddHotel(store *db.Store, name, location string, rating int) *types.Hotel {
	hotel := &types.Hotel{
		Name:     name,
		Location: location,
		Rating:   rating,
		Rooms:    []primitive.ObjectID{},
	}
	insertedHotel, err := store.Hotel.CreateHotel(context.TODO(), hotel)

	room_1 := AddRoom(store, "small", true, 100, insertedHotel.ID)
	room_2 := AddRoom(store, "medium", false, 120, insertedHotel.ID)
	room_3 := AddRoom(store, "large", true, 120, insertedHotel.ID)

	rooms := []*types.Room{room_1, room_2, room_3}
	for _, room := range rooms {
		insertedHotel.Rooms = append(insertedHotel.Rooms, room.ID)
	}

	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}

func AddRoom(store *db.Store, size string, seaside bool, price float64, hotelID primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: seaside,
		Price:   price,
		HotelID: hotelID,
	}
	insertedRoom, err := store.Room.CreateRoom(context.TODO(), room)
	if err != nil {
		log.Fatal(err)
	}

	return insertedRoom
}

func AddBooking(store *db.Store, userID, roomID primitive.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := store.Booking.BookRoom(context.TODO(), booking)
	if err != nil {
		log.Fatal(err)
	}

	return insertedBooking
}
