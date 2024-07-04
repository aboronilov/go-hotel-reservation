package main

import (
	"flag"

	"github.com/aboronilov/go-hotel-reservation/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	listenAddr := flag.String("listenAddr", ":5000", "The listen address of the server")
	flag.Parse()
	app := fiber.New()
	apiv1 := app.Group("/api/v1")
	apiv1.Get("/user", api.HandleListUsers)
	apiv1.Get("/user/:id", api.HandleListUser)
	app.Listen(*listenAddr)
}
