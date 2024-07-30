package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aboronilov/go-hotel-reservation/api/middleware"
	"github.com/aboronilov/go-hotel-reservation/db/fixtures"
	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestAdminCanGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		user           = fixtures.AddUser(db.store, "john", "smith", false)
		adminUser      = fixtures.AddUser(db.store, "james", "bond", true)
		hotel          = fixtures.AddHotel(db.store, "ibis", "paris", 5)
		from           = time.Now().AddDate(0, 0, 1)
		till           = time.Now().AddDate(0, 0, 6)
		booking        = fixtures.AddBooking(db.store, user.ID, hotel.Rooms[0], from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.store.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.store)
	)

	admin.Get("/", bookingHandler.HandleListBookings)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", CreateTokenFromUser(adminUser))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	var bookings []*types.Booking
	if err = json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}

	if len(bookings) != 1 || bookings[0].ID != booking.ID {
		t.Fatalf("expected booking to be returned, got %+v", bookings)
	}

	if bookings[0].UserID != booking.UserID {
		t.Fatalf("expected user ID to match, got %s, expected %s", bookings[0].UserID, booking.UserID)
	}

	if bookings[0].RoomID != booking.RoomID {
		t.Fatalf("expected rooom ID to match, got %s, expected %s", bookings[0].RoomID, booking.RoomID)
	}

	if bookings[0].NumPersons != booking.NumPersons {
		t.Fatalf("expected number of persons to match, got %d, expected %d", bookings[0].NumPersons, booking.NumPersons)
	}
}

func TestUserCannotGetBookingsOfOtherUsers(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		user           = fixtures.AddUser(db.store, "john", "smith", false)
		hotel          = fixtures.AddHotel(db.store, "ibis", "paris", 5)
		from           = time.Now().AddDate(0, 0, 1)
		till           = time.Now().AddDate(0, 0, 6)
		booking        = fixtures.AddBooking(db.store, user.ID, hotel.Rooms[0], from, till)
		app            = fiber.New()
		admin          = app.Group("/", middleware.JWTAuthentication(db.store.User), middleware.AdminAuth)
		bookingHandler = NewBookingHandler(db.store)
	)

	_ = booking

	admin.Get("/", bookingHandler.HandleListBookings)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status code 401, got %d", resp.StatusCode)
	}
}

func TestUserCanGetOwnBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		user           = fixtures.AddUser(db.store, "john", "smith", false)
		hotel          = fixtures.AddHotel(db.store, "ibis", "paris", 5)
		from           = time.Now().AddDate(0, 0, 1)
		till           = time.Now().AddDate(0, 0, 6)
		booking        = fixtures.AddBooking(db.store, user.ID, hotel.Rooms[0], from, till)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(db.store.User))
		bookingHandler = NewBookingHandler(db.store)
	)

	route.Get("/:id", bookingHandler.HandleRetrieveBooking)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Add("Authorization", CreateTokenFromUser(user))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	var actualBooking *types.Booking
	if err = json.NewDecoder(resp.Body).Decode(&actualBooking); err != nil {
		t.Fatal(err)
	}

	if actualBooking.ID != booking.ID {
		t.Fatalf("expected booking to be returned, got %+v", actualBooking)
	}

	if actualBooking.UserID != booking.UserID {
		t.Fatalf("expected user ID to match, got %s, expected %s", actualBooking.UserID, booking.UserID)
	}

	if actualBooking.RoomID != booking.RoomID {
		t.Fatalf("expected rooom ID to match, got %s, expected %s", actualBooking.RoomID, booking.RoomID)
	}

	if actualBooking.NumPersons != booking.NumPersons {
		t.Fatalf("expected number of persons to match, got %d, expected %d", actualBooking.NumPersons, booking.NumPersons)
	}
}

func TestUserWithoutAuthCantGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	var (
		user           = fixtures.AddUser(db.store, "john", "smith", false)
		hotel          = fixtures.AddHotel(db.store, "ibis", "paris", 5)
		from           = time.Now().AddDate(0, 0, 1)
		till           = time.Now().AddDate(0, 0, 6)
		booking        = fixtures.AddBooking(db.store, user.ID, hotel.Rooms[0], from, till)
		app            = fiber.New()
		route          = app.Group("/", middleware.JWTAuthentication(db.store.User))
		bookingHandler = NewBookingHandler(db.store)
	)

	route.Get("/:id", bookingHandler.HandleRetrieveBooking)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
}
