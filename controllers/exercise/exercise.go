package exercise

import (
	"phuhung273/progress-tracking/db"
	"phuhung273/progress-tracking/middleware"
	"phuhung273/progress-tracking/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

func index(c *fiber.Ctx) error {
	sess, _ := middleware.SessionStore.Get(c)

	userId := sess.Get("user_id")

	var items []models.Exercise
	db.DB.Where("user_id = ?", userId).Order("id DESC").Joins("Category").Joins("SecondaryCategory").Limit(10).Find(&items)

	return c.Render("exercise/index.html", fiber.Map{
		"title": "Exercise",
		"items": items,
	})
}

func create(c *fiber.Ctx) error {
	var allCategories []models.Setting
	db.DB.Where("type IN ?", []string{"CATEGORY", "SECONDARY_CATEGORY"}).Find(&allCategories)

	categories := []models.Setting{}
	secondaryCategories := []models.Setting{}
	for _, setting := range allCategories {
		if setting.Type == "CATEGORY" {
			categories = append(categories, setting)
		} else if setting.Type == "SECONDARY_CATEGORY" {
			secondaryCategories = append(secondaryCategories, setting)
		}
	}

	return c.Render("exercise/form.html", fiber.Map{
		"title": "Exercise",
		"categories": categories,
		"secondaryCategories": secondaryCategories,
	})
}

func store(c *fiber.Ctx) error {
	sess, _ := middleware.SessionStore.Get(c)

	cate, _ := strconv.Atoi(c.FormValue("category"))
	secondaryCate, _ := strconv.Atoi(c.FormValue("secondary_category"))
	userId := sess.Get("user_id").(uint)

	item := &models.Exercise{ 
		CategoryID: uint(cate),
		UserID: uint(userId),
		Result: datatypes.JSON([]byte(`{"test": "test"}`)),
	}

	if secondaryCate > 0 {
		v := uint(secondaryCate)
		item.SecondaryCategoryID = &v
	}

	db.DB.Create(item)
	return c.Redirect("/exercise")
}

func edit(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Exercise
	exist := db.DB.First(&item, id)
	if exist.RowsAffected == 0 {
		return c.Redirect( "/setting")
	}

	var allCategories []models.Setting
	db.DB.Where("type IN ?", []string{"CATEGORY", "SECONDARY_CATEGORY"}).Find(&allCategories)

	categories := []models.Setting{}
	secondaryCategories := []models.Setting{}
	for _, setting := range allCategories {
		if setting.Type == "CATEGORY" {
			categories = append(categories, setting)
		} else if setting.Type == "SECONDARY_CATEGORY" {
			secondaryCategories = append(secondaryCategories, setting)
		}
	}

	return c.Render("exercise/form.html", fiber.Map{
		"title": "Setting",
		"item": item,
		"categories": categories,
		"secondaryCategories": secondaryCategories,
	})
}

func update(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Exercise
	exist := db.DB.First(&item, id)
	if exist.RowsAffected == 0 {
		return c.Redirect("/exercise")
	}

	cate, _ := strconv.Atoi(c.FormValue("category"))
	secondaryCate, _ := strconv.Atoi(c.FormValue("secondary_category"))

	item.CategoryID = uint(cate)
	item.Result = datatypes.JSON([]byte(`{"test": "test"}`))

	if secondaryCate > 0 {
		v := uint(secondaryCate)
		item.SecondaryCategoryID = &v
	}

	db.DB.Save(&item)
	return c.Redirect("/exercise")
}

func delete(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Exercise
	exist := db.DB.First(&item, id)
	if exist.RowsAffected == 0 {
		return c.Redirect("/exercise")
	}

	db.DB.Delete(&models.Exercise{}, id)
	return c.Redirect( "/exercise")
}

func Routing(router *fiber.App) {
	router.Route("/exercise", func(router fiber.Router) {
		router.Use(middleware.Auth)
		
		router.Get("/", index)
		router.Get("/create", create)
		router.Post("/", store)
		router.Get("/edit/:id", edit)
		router.Post("/edit/:id", update)
		router.Get("/delete/:id", delete)
	})
}