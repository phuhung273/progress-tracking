package exercise

import (
	"phuhung273/progress-tracking/controllers"
	"phuhung273/progress-tracking/db"
	"phuhung273/progress-tracking/middleware"
	"phuhung273/progress-tracking/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func index(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(uint)

	var count int64
	countQueryChannel := make(chan bool)
	go func() {
		db.DB.Model(&models.Exercise{}).Where("user_id = ?", userId).Count(&count)
		countQueryChannel <- true
	}()

	var items []models.Exercise
	db.DB.Where("user_id = ?", userId).Order("id DESC").Joins("Category").Joins("SecondaryCategory").Preload("Results.Criteria").Limit(10).Find(&items)

	<- countQueryChannel

	return controllers.ListResponse(c, items, count)
}

func store(c *fiber.Ctx) error {
	var item models.Exercise

	c.BodyParser(&item)
	userId := c.Locals("user_id").(uint)

	item.UserID = userId
	if *item.SecondaryCategoryID == 0 {
		item.SecondaryCategoryID = nil
	}
	
	db.DB.Create(&item)

	return c.JSON(item)
}

func edit(c *fiber.Ctx) error {
	id := c.Params("id")

	categories := []models.Setting{}
	secondaryCategories := []models.Setting{}
	criterias := []models.Setting{}
	querySettingChannel := make(chan bool)

	go func() {
		settings, _ := db.DB.Model(&models.Setting{}).Rows()
		defer settings.Close()

		for settings.Next() {
			
			var setting models.Setting
			db.DB.ScanRows(settings, &setting)

			if setting.Type == "CATEGORY" {
				categories = append(categories, setting)
			} else if setting.Type == "SECONDARY_CATEGORY" {
				secondaryCategories = append(secondaryCategories, setting)
			} else if setting.Type == "CRITERIA" {
				criterias = append(criterias, setting)
			}
		}
		querySettingChannel <- true
	}()

	var item models.Exercise
	exist := db.DB.Preload("Results").First(&item, id)
	if exist.RowsAffected == 0 {
		return c.Redirect( "/setting")
	}

	<- querySettingChannel

	return c.Render("exercise/form.html", fiber.Map{
		"title": "Setting",
		"item": item,
		"secondaryCategory": *item.SecondaryCategoryID,
		"categories": categories,
		"secondaryCategories": secondaryCategories,
		"criterias": criterias,
	})
}

func update(c *fiber.Ctx) error {
	id := c.Params("id")

	deleteResultChannel := make(chan bool)
	go func() {
		db.DB.Where("exercise_id = ?", id).Delete(&models.Result{})
		deleteResultChannel <- true
	}()

	var item models.Exercise
	exist := db.DB.First(&item, id)
	if exist.RowsAffected == 0 {
		return c.Redirect("/exercise")
	}

	form, _ :=c.MultipartForm()

	cate, _ := strconv.Atoi(form.Value["category"][0])
	secondaryCate, _ := strconv.Atoi(form.Value["secondary_category"][0])
	criterias := form.Value["criteria"]
	values := form.Value["value"]

	item.CategoryID = uint(cate)

	if secondaryCate > 0 {
		v := uint(secondaryCate)
		item.SecondaryCategoryID = &v
	}

	db.DB.Save(&item)
	
	results := []models.Result{}
	for i := 0; i < len(criterias); i++ {
		criteria, _ := strconv.Atoi(criterias[i])
		value, _ := strconv.Atoi(values[i])
		results = append(results, models.Result{
			Value: uint(value),
			CriteriaID: uint(criteria),
			ExerciseID: item.ID,
		})
	}
	db.DB.Create(&results)

	return c.Redirect("/exercise")
}

func delete(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Exercise
	exist := db.DB.First(&item, id)
	if exist.RowsAffected == 0 {
		return c.SendStatus(404)
	}

	db.DB.Delete(&models.Exercise{}, id)
	return c.SendStatus(200)
}

func Routing(router *fiber.App) {
	router.Route("/exercise", func(router fiber.Router) {
		router.Use(middleware.Auth)
		
		router.Get("/", index)
		router.Post("/", store)
		router.Get("/:id", edit)
		router.Put("/:id", update)
		router.Delete("/:id", delete)
	})
}