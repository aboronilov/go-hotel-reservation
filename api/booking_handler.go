package api

import (
	"errors"

	"github.com/aboronilov/go-hotel-reservation/db"
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
		return ErrorNotFound()
	}

	return c.JSON(bookings)
}

func (h *BookinHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrorInvalidID()
	}

	booking, err := h.store.Booking.GetBookingByID(c.Context(), oid)
	if err != nil {
		return ErrorNotFound()
	}

	user, err := getAuthUser(c)
	if err != nil || booking.UserID != user.ID {
		return ErrorUnauthorized()
	}

	err = h.store.Booking.UpdateBooking(c.Context(), oid, bson.M{"canceled": true})
	if err != nil {
		return ErrorBadRequest()
	}

	return c.JSON(map[string]string{"message": "Booking canceled"})
}

// only owner
func (h *BookinHandler) HandleRetrieveBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrorInvalidID()
	}

	booking, err := h.store.Booking.GetBookingByID(c.Context(), oid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrorNotFound()
		}
		return err
	}

	user, err := getAuthUser(c)
	if err != nil || booking.UserID != user.ID {
		return ErrorUnauthorized()
	}

	return c.JSON(booking)
}
