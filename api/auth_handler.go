package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aboronilov/go-hotel-reservation/db"
	"github.com/aboronilov/go-hotel-reservation/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  *types.User `json:"user"`
}

type genericResponse struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResponse{
		Message: "Invalid credentials",
		Type:    "error",
	})
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	// fmt.Println("Auth Params: ", authParams)

	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		fmt.Println("err: ", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	if !types.IsValidPassword(user.HashedPassword, authParams.Password) {
		return invalidCredentials(c)
	}

	response := AuthResponse{
		Token: CreateTokenFromUser(user),
		User:  user,
	}

	return c.JSON(response)
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 72).Unix()
	claims := jwt.MapClaims{
		"id":      user.ID,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println("Failed signing token")
	}

	return tokenStr
}
