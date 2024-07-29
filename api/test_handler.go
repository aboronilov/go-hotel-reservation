package api

import (
	"context"
	"log"
	"testing"

	"github.com/aboronilov/go-hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	client *mongo.Client
	store  *db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.TestDBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)

	return &testdb{
		client: client,
		store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Hotel:   hotelStore,
			Booking: db.NewMongoBookingStore(client),
		},
	}
}
