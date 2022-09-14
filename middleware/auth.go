package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var SessionStore *session.Store

func Auth(c *fiber.Ctx) error {
	sess, _ := SessionStore.Get(c)
	user := sess.Get("user_id")
	if user == nil {
		return c.Redirect("/auth/login")
	}
	return c.Next()
}