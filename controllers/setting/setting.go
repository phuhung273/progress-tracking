package setting

import (
	"phuhung273/progress-tracking/db"
	"phuhung273/progress-tracking/middleware"
	"phuhung273/progress-tracking/models"

	"github.com/gofiber/fiber/v2"
)

func index(c *fiber.Ctx) error {
	var items []models.Setting
	db.DB.Order("id DESC").Limit(10).Find(&items)

	return c.Render("setting/index.html", fiber.Map{
		"title": "Setting",
		"items": items,
	})
}

func create(c *fiber.Ctx) error {
	return c.Render("setting/form.html", fiber.Map{
		"title": "Setting",
	})
}

func store(c *fiber.Ctx) error {
	typeValue := c.FormValue("type")
	name := c.FormValue("name")

	db.DB.Create(&models.Setting{ Type: typeValue, Name: name })
	return c.Redirect("/settings")
}

func edit(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Setting
	exist := db.DB.First(&item, id)
	if exist.RowsAffected == 0 {
		return c.Redirect( "/settings")
	}

	return c.Render("setting/form.html", fiber.Map{
		"title": "Setting",
		"item": item,
	})
}

func update(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Setting
	exist := db.DB.First(&item, id)
	if exist.RowsAffected == 0 {
		return c.Redirect("/settings")
	}

	item.Name = c.FormValue("name")
	item.Type = c.FormValue("type")
	db.DB.Save(&item)
	return c.Redirect("/settings")
}

func delete(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Setting
	exist := db.DB.First(&item, id)
	if exist.RowsAffected == 0 {
		return c.Redirect("/settings")
	}

	db.DB.Delete(&models.Setting{}, id)
	return c.Redirect( "/settings")
}

func Routing(router *fiber.App) {
	router.Route("/settings", func(router fiber.Router) {
		router.Use(middleware.Auth)
		
		router.Get("/", index)
		router.Get("/create", create)
		router.Post("/", store)
		router.Get("/edit/:id", edit)
		router.Post("/edit/:id", update)
		router.Get("/delete/:id", delete)
	})
}