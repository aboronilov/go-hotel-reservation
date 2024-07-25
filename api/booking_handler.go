package api

import (
	"errors"

	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookinHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookinHandler {
	return &BookinHandler{
		store: store,
	}
}

// admin auth
func (h *BookinHandler) HandleListBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(bookings)
}

// only owner
func (h *BookinHandler) HandleRetrieveBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	booking, err := h.store.Booking.GetBookingByID(c.Context(), oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok || booking.UserID != user.ID {
		return c.Status(fiber.StatusUnauthorized).JSON(map[string]string{"error": "Unauthorized"})
	}

	return c.JSON(booking)
}
