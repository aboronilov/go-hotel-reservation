package api

import (
	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return nil, c.Status(fiber.StatusUnauthorized).JSON(map[string]string{"error": "Unauthorized"})
	}

	return user, nil
}
