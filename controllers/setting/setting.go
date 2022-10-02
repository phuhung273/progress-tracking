package setting

import (
	"phuhung273/progress-tracking/controllers"
	"phuhung273/progress-tracking/db"
	"phuhung273/progress-tracking/middleware"
	"phuhung273/progress-tracking/models"

	"github.com/gofiber/fiber/v2"
)

func index(c *fiber.Ctx) error {
	var count int64
	countQueryChannel := make(chan bool)
	go func() {
		db.DB.Model(&models.Setting{}).Count(&count)
		countQueryChannel <- true
	}()

	var items []models.Setting
	db.DB.Order("id DESC").Limit(10).Find(&items)

	<- countQueryChannel

	return controllers.ListResponse(c, items, count)
}

func store(c *fiber.Ctx) error {
	item := new(models.Setting)
	c.BodyParser(item)

	db.DB.Create(&item)
	return c.JSON(item)
}

func edit(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Setting
	exist := db.DB.First(&item, id)
	if exist.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	return c.JSON(item)
}

func update(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Setting
	exist := db.DB.First(&item, id)
	if exist.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	newItem := new(models.Setting)
	c.BodyParser(newItem)

	item.Name = newItem.Name
	item.Type = newItem.Type

	db.DB.Save(&item)
	return c.JSON(item)
}

func delete(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Setting
	exist := db.DB.First(&item, id)
	if exist.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	db.DB.Delete(&models.Setting{}, id)
	return c.SendStatus(200)
}

func Routing(router *fiber.App) {
	router.Route("/setting", func(router fiber.Router) {
		router.Use(middleware.Auth)
		
		router.Get("/", index)
		router.Post("/", store)
		router.Get("/:id", edit)
		router.Put("/:id", update)
		router.Delete("/:id", delete)
	})
}