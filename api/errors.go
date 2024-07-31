package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func ErrorHandler(c *fiber.Ctx, err error) error {
	if apiError, ok := err.(*Error); ok {
		return c.Status(apiError.Code).JSON(map[string]string{"error": apiError.Message})
	}
	newError := NewError(http.StatusInternalServerError, err.Error())
	return c.Status(newError.Code).JSON(map[string]string{"error": newError.Message})
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func (e *Error) Error() string {
	return e.Message
}

func ErrorInvalidID() *Error {
	return NewError(http.StatusBadRequest, "invalid id")
}

func ErrorUnauthorized() *Error {
	return NewError(http.StatusUnauthorized, "unauthorized")
}

func ErrorNotFound() *Error {
	return NewError(http.StatusNotFound, "not found")
}

func ErrorBadRequest() *Error {
	return NewError(http.StatusBadRequest, "invalid JSON request")
}
