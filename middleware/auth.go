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
	// sess, _ := SessionStore.Get(c)
	// user := sess.Get("user_id")
	// if user == nil {
	// 	return c.Redirect("/auth/login")
	// }

	authorization := c.Get("Authorization")
	bearerToken := strings.Split(authorization, " ")[1]

	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
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
		c.Locals("user_id", claims["user_id"])
	} else {
		return c.Status(401).JSON(fiber.Map{
			"message": "Token invalid",
		})
	}
	
	return c.Next()
}