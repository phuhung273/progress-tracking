package home

import (
	"phuhung273/progress-tracking/middleware"

	"github.com/gofiber/fiber/v2"
)

func index(c *fiber.Ctx) error {
	return c.Render("home.html", fiber.Map{
		"title": "Home",
	})
}

func Routing(router *fiber.App) {
	router.Get("/dashboard", middleware.Auth, index)
}