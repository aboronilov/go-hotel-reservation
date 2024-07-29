package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/aboronilov/go-hotel-reservation/db/fixtures"
)

func TestAdminGetBookings(t *testing.T) {
	db := setup(t)
	defer db.teardown(t)

	user := fixtures.AddUser(db.store, "john", "smith", false)
	hotel := fixtures.AddHotel(db.store, "ibis", "paris", 5)

	from := time.Now().AddDate(0, 0, 1)
	till := time.Now().AddDate(0, 0, 6)
	booking := fixtures.AddBooking(db.store, user.ID, hotel.Rooms[0], from, till)

	fmt.Println(booking)
}
