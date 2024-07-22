package db

const (
	TestDBNAME         = "test-hotel-reservation"
	DBNAME             = "hotel-reservation"
	DBURI              = "mongodb://localhost:27017"
	HOTEL_COLLECTION   = "hotels"
	USERS_COLLECTION   = "users"
	ROOM_COLLECTION    = "rooms"
	BOOKING_COLLECTION = "bookings"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
