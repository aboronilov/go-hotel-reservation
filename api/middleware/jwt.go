package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func JWTAuthentication(c *fiber.Ctx) error {
	fmt.Println("--- JWT Auth")

	headers := c.GetReqHeaders()
	token, ok := headers["Authorization"]
	if !ok || len(token) != 1 {
		return fmt.Errorf("Unauthorized")
	}

	if err := parseToken(token[0]); err != nil {
		return err
	}

	return nil
}

func parseToken(tokenStr string) error {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("Invalid signing method: ", token.Header["alg"])
			return nil, fmt.Errorf("Unauthorized")
		}

		secret := os.Getenv("JWT_SECRET")
		fmt.Println("secret: ", secret)
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("Failed to parse JWT: ", err)
		return fmt.Errorf("Unauthorized")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims)
	}

	return fmt.Errorf("Unauthorized")

}
