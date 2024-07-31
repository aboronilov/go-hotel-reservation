package api

import (
	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok || !user.IsAdmin {
		return ErrorUnauthorized()
	}

	return c.Next()
}
