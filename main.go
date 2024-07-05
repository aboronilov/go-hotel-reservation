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

const (
	uri = "mongodb://localhost:27017"
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

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	// userJames := &types.User{
	// 	FirstName: "James",
	// 	LastName:  "Smith",
	// }
	// doc, err := client.Database(db.DBNAME).Collection("users").InsertOne(ctx, userJames)
	// if err != nil {
	// 	log.Fatal("Error inserting user", err)
	// }
	// fmt.Println(doc)

	// cur, err := client.Database(db.DBNAME).ListCollections(ctx, bson.M{})
	// if err != nil {
	// 	log.Fatal("Error finding collections", err)
	// }

	// var collections []*mongo.Collection
	// for cur.Next(ctx) {
	// 	var collection mongo.Collection
	// 	err = cur.Decode(&collections)
	// 	if err != nil {
	// 		log.Fatal("Error decoding collection", err)
	// 	}
	// 	collections = append(collections, &collection)
	// }
	// fmt.Println(collections)

	app := fiber.New(config)
	apiv1 := app.Group("/api/v1")

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	apiv1.Get("/user", userHandler.HandleListUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	app.Listen(*listenAddr)
}
