package api

import (
	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func HandleListUsers(c *fiber.Ctx) error {
	u := types.User{
		ID:        "1",
		FirstName: "James",
		LastName:  "Smith",
	}
	return c.JSON(u)
}

func HandleListUser(c *fiber.Ctx) error {
	return c.JSON(map[string]string{"User": "James"})
}
