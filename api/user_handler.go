package api

import (
	"context"

	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandleGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)

	user, err := h.userStore.GetUserByID(ctx, id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleListUsers(c *fiber.Ctx) error {
	u := types.User{
		ID:        "1",
		FirstName: "James",
		LastName:  "Smith",
	}
	return c.JSON(u)
}
