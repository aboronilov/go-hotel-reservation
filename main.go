package main

import (
	"context"
	"flag"
	"log"

	"github.com/aboronilov/go-hotel-reservation/api"
	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the server")
	flag.Parse()

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	// stores
	hotelStore := db.NewMongoHotelStore(client, false)
	roomStore := db.NewMongoRoomStore(client, hotelStore, false)
	userStore := db.NewMongoUserStore(client, false)
	bookingStore := db.NewMongoBookingStore(client, false)
	store := &db.Store{
		User:    userStore,
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: bookingStore,
	}

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1", api.JWTAuthentication(userStore))
	auth := app.Group("/api")
	admin := apiv1.Group("/admin", api.AdminAuth)

	// user
	userHandler := api.NewUserHandler(userStore)
	apiv1.Get("/user", userHandler.HandleListUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Post("/user", userHandler.HandleCreateUser)
	apiv1.Put("/user/:id", userHandler.HandleUpdateUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// auth
	authHandler := api.NewAuthHandler(userStore)
	auth.Post("/auth", authHandler.HandleAuthenticate)

	// room
	roomHandler := api.NewRoomHandler(store)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	apiv1.Get("/room", roomHandler.HandleListRooms)

	// bookings
	bookingHandler := api.NewBookingHandler(store)
	apiv1.Get("/booking/:id", bookingHandler.HandleRetrieveBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// admin
	admin.Get("/booking", bookingHandler.HandleListBookings)

	// hotel
	hotelHandler := api.NewHotelHandler(store)
	apiv1.Get("/hotel", hotelHandler.HandleListHotels)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiv1.Get("/hotel/:id", hotelHandler.HandleRetrieveHotel)

	app.Listen(*listenAddr)
}
