package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/golang-jwt/jwt/v4"
)

var SessionStore *session.Store

func Auth(c *fiber.Ctx) error {

	authorization := c.Get("Authorization")
	bearerToken := strings.Split(authorization, " ")[1]

	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		v, _ := err.(*jwt.ValidationError)
		if v.Errors == jwt.ValidationErrorExpired {
			return c.Status(401).JSON(fiber.Map{
				"message": "Token expired",
			})
		}
		return c.Status(401).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.Locals("user_id", uint(claims["user_id"].(float64)))
	} else {
		return c.Status(401).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}
	
	return c.Next()
}