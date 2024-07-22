package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RoomHandler struct {
	store *db.Store
}

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p *BookRoomParams) validate() error {
	now := time.Now()
	if p.FromDate.Before(now) || p.TillDate.Before(now) {
		return fiber.NewError(http.StatusBadRequest, "Invalid dates: fromDate and tillDate should be in the future")
	}

	if p.FromDate.After(p.TillDate) {
		return fiber.NewError(http.StatusBadRequest, "Invalid dates: fromDate should be before tillDate")
	}

	if p.NumPersons <= 0 {
		return fiber.NewError(http.StatusBadRequest, "Invalid number of persons")
	}

	return nil
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HandleListRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}

	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(genericResponse{
			Message: "Internal server error",
			Type:    "Error",
		})
	}

	ok, err = h.isRoomAvailiable(c.Context(), roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResponse{
			Message: fmt.Sprintf("Room %s is already booked between %s and %s", roomID.Hex(), params.FromDate.Format(time.RFC3339), params.TillDate.Format(time.RFC3339)),
			Type:    "Error",
		})
	}

	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	inserted, err := h.store.Booking.BookRoom(c.Context(), &booking)
	if err != nil {
		return err
	}

	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailiable(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	filter := bson.M{
		"roomID":   roomID,
		"fromDate": bson.M{"$lte": params.TillDate},
		"tillDate": bson.M{"$gte": params.FromDate},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, filter)

	if err != nil {
		return false, err
	}

	ok := len(bookings) == 0
	return ok, nil
}
