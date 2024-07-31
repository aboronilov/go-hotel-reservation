package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
)

func TestCreateUser(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New(fiber.Config{ErrorHandler: ErrorHandler})
	userHandler := NewUserHandler(tdb.store.User)
	app.Post("/", userHandler.HandleCreateUser)

	params := types.CreateUserParams{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
	}
	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	_ = json.NewDecoder(resp.Body).Decode(&user)

	if user.Email != params.Email {
		t.Error("expected email to match", params.Email, "!=", user.Email)
	}
	if user.FirstName != params.FirstName {
		t.Error("expected first name to match", params.FirstName, "!=", user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Error("expected last name to match", params.LastName, "!=", user.LastName)
	}
	if len(user.ID) == 0 {
		t.Errorf("expected user ID to be set")
	}
}
