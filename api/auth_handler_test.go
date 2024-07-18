package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func createTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "James",
		LastName:  "Bond",
		Email:     "james@bond.com",
		Password:  "secret",
	})

	if err != nil {
		t.Fatal(err)
	}

	newUser, err := userStore.CreateUser(context.Background(), user)
	if err != nil {
		t.Fatal(err)
	}

	return newUser
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := createTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "james@bond.com",
		Password: "secret",
	}
	b, _ := json.Marshal(authParams)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}

	var authResponse AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		t.Fatal(err)
	}

	if authResponse.Token == "" {
		t.Fatalf("expected token to be set")
	}
	// Hashed Password is empty because we don't return in JSON response
	insertedUser.HashedPassword = ""
	if !reflect.DeepEqual(insertedUser, authResponse.User) {
		t.Fatalf("expected user to match: %+v, got %+v", insertedUser, authResponse.User)
	}
}

func TestAuthenticateFailsWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	createTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "james@bond.com",
		Password: "false_password",
	}
	b, _ := json.Marshal(authParams)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code 400, got %d", resp.StatusCode)
	}

	var genericResponse genericResponse
	err = json.NewDecoder(resp.Body).Decode(&genericResponse)
	if err != nil {
		t.Fatal(err)
	}
	if genericResponse.Message != "Invalid credentials" {
		t.Fatalf("expected message 'Invalid credentials', got '%s'", genericResponse.Message)
	}
	if genericResponse.Type != "error" {
		t.Fatalf("expected type 'error', got '%s'", genericResponse.Type)
	}
}

func TestAuthenticateFailsWithWrongEmail(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	createTestUser(t, tdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "bond@james-bond.com",
		Password: "secret",
	}
	b, _ := json.Marshal(authParams)

	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code 400, got %d", resp.StatusCode)
	}

	var genericResponse genericResponse
	err = json.NewDecoder(resp.Body).Decode(&genericResponse)
	if err != nil {
		t.Fatal(err)
	}
	if genericResponse.Message != "Invalid credentials" {
		t.Fatalf("expected message 'Invalid credentials', got '%s'", genericResponse.Message)
	}
	if genericResponse.Type != "error" {
		t.Fatalf("expected type 'error', got '%s'", genericResponse.Type)
	}
}
