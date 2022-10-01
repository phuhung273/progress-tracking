package router

import (
	"phuhung273/progress-tracking/controllers/exercise"
	"phuhung273/progress-tracking/controllers/home"
	"phuhung273/progress-tracking/controllers/setting"
	"phuhung273/progress-tracking/controllers/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var Router *fiber.App

func Route() {
	Router.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/dashboard")
	})

	user.Routing(Router)
	setting.Routing(Router)
	exercise.Routing(Router)
	home.Routing(Router)
}
func Init() {
	
	// engine := html.New("views", ".html")
	// engine.Reload(true)

	Router = fiber.New(fiber.Config{
		// Views: engine,
	})

	Router.Use(cors.New())

	// middleware.SessionStore = session.New()
}