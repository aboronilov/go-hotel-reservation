package main

import (
	"context"
	"flag"
	"log"

	"github.com/aboronilov/go-hotel-reservation/api"
	"github.com/aboronilov/go-hotel-reservation/api/middleware"
	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: func(c *fiber.Ctx, err error) error {
		return c.JSON(map[string]string{"error": err.Error()})
	},
}

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the server")
	flag.Parse()

	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1", middleware.JWTAuthentication)
	auth := app.Group("/api")

	// stores
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	userStore := db.NewMongoUserStore(client)
	store := &db.Store{
		User:  userStore,
		Hotel: hotelStore,
		Room:  roomStore,
	}

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
	// roomHandler := api.NewRoomHandler(roomStore)

	// hotel
	hotelHandler := api.NewHotelHandler(store)
	apiv1.Get("/hotel", hotelHandler.HandleListHotels)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	apiv1.Get("/hotel/:id", hotelHandler.HandleRetrieveHotel)
	// apiv1.Post("/hotel", hotelHandler.HandleCreateHotel)
	// apiv1.Put("/hotel/:id", hotelHandler.HandleUpdateHotel)
	// apiv1.Delete("/hotel/:id", hotelHandler.HandleDeleteHotel)

	app.Listen(*listenAddr)
}
