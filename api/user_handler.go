package api

import (
	"errors"
	"fmt"

	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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
	id := c.Params("id")

	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrorNotFound()
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleListUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return ErrorNotFound()
	}

	return c.JSON(users)
}

func (h *UserHandler) HandleCreateUser(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}

	if errors := params.Validate(); len(errors) != 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}

	createdUser, err := h.userStore.CreateUser(c.Context(), user)
	if err != nil {
		return err
	}

	return c.JSON(createdUser)
}

func (h *UserHandler) HandleUpdateUser(c *fiber.Ctx) error {
	var (
		// update bson.M
		params types.UpdateUserParams
		userId = c.Params("id")
	)

	id, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return ErrorInvalidID()
	}

	if err = c.BodyParser(&params); err != nil {
		return ErrorBadRequest()
	}

	filter := bson.M{"_id": id}

	err = h.userStore.UpdateUserByID(c.Context(), filter, params)
	if err != nil {
		return ErrorBadRequest()
	}

	return c.JSON(map[string]string{"msg": fmt.Sprintf("user %s updated", userId)})
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	if err := h.userStore.DeleteUserByID(c.Context(), userId); err != nil {
		return err
	}

	return c.JSON(map[string]string{"msg": fmt.Sprintf("user %s deleted", userId)})
}
