package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/aboronilov/go-hotel-reservation/db/fixtures"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)
	insertedUser := fixtures.AddUser(tdb.store, "jason_1", "bourne", true)

	// fmt.Println("insertedUser --->", insertedUser)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "jason_1_bourne@ctu.com",
		Password: "jason_1_bourne",
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
	// createTestUser(t, tdb.store.User)
	fixtures.AddUser(tdb.store, "james", "bond", true)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	authParams := AuthParams{
		Email:    "james_bond@ctu.com",
		Password: "wrong_password",
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
	fixtures.AddUser(tdb.store, "james", "bond", true)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.store.User)
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
