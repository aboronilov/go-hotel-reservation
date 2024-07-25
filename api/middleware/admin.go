package middleware

import (
	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok || !user.IsAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(map[string]string{"error": "Unauthorized"})
	}

	return c.Next()
}
