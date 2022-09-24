package exercise

import (
	"phuhung273/progress-tracking/db"
	"phuhung273/progress-tracking/middleware"
	"phuhung273/progress-tracking/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func index(c *fiber.Ctx) error {
	sess, _ := middleware.SessionStore.Get(c)

	userId := sess.Get("user_id")

	var items []models.Exercise
	db.DB.Where("user_id = ?", userId).Order("id DESC").Joins("Category").Joins("SecondaryCategory").Preload("Results.Criteria").Limit(10).Find(&items)

	return c.Render("exercise/index.html", fiber.Map{
		"title": "Exercise",
		"items": items,
	})
}

func create(c *fiber.Ctx) error {
	var settings []models.Setting
	db.DB.Find(&settings)

	categories := []models.Setting{}
	secondaryCategories := []models.Setting{}
	criterias := []models.Setting{}
	for _, setting := range settings {
		if setting.Type == "CATEGORY" {
			categories = append(categories, setting)
		} else if setting.Type == "SECONDARY_CATEGORY" {
			secondaryCategories = append(secondaryCategories, setting)
		} else if setting.Type == "CRITERIA" {
			criterias = append(criterias, setting)
		}
	}

	return c.Render("exercise/form.html", fiber.Map{
		"title": "Exercise",
		"categories": categories,
		"secondaryCategories": secondaryCategories,
		"criterias": criterias,
	})
}

func store(c *fiber.Ctx) error {
	sess, _ := middleware.SessionStore.Get(c)

	form, _ :=c.MultipartForm()

	cate, _ := strconv.Atoi(form.Value["category"][0])
	secondaryCate, _ := strconv.Atoi(form.Value["secondary_category"][0])
	criterias := form.Value["criteria"]
	values := form.Value["value"]
	userId := sess.Get("user_id").(uint)

	item := models.Exercise{ 
		CategoryID: uint(cate),
		UserID: uint(userId),
	}

	if secondaryCate > 0 {
		v := uint(secondaryCate)
		item.SecondaryCategoryID = &v
	}
	
	db.DB.Create(&item)

	for i := 0; i < len(criterias); i++ {
		criteria, _ := strconv.Atoi(criterias[i])
		value, _ := strconv.Atoi(values[i])
		result := models.Result{
			Value: uint(value),
			CriteriaID: uint(criteria),
			ExerciseID: item.ID,
		}
		db.DB.Create(&result)
	}

	return c.Redirect("/exercise")
}

func edit(c *fiber.Ctx) error {
	id := c.Params("id")

	var item models.Exercise
	exist := db.DB.Preload("Results").First(&item, id)
	if exist.RowsAffected == 0 {
		return c.Redirect( "/setting")
	}

	var settings []models.Setting
	db.DB.Find(&settings)

	categories := []models.Setting{}
	secondaryCategories := []models.Setting{}
	criterias := []models.Setting{}
	for _, setting := range settings {
		if setting.Type == "CATEGORY" {
			categories = append(categories, setting)
		} else if setting.Type == "SECONDARY_CATEGORY" {
			secondaryCategories = append(secondaryCategories, setting)
		} else if setting.Type == "CRITERIA" {
			criterias = append(criterias, setting)
		}
	}

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

	db.DB.Where("exercise_id = ?", item.ID).Delete(&models.Result{})
	
	for i := 0; i < len(criterias); i++ {
		criteria, _ := strconv.Atoi(criterias[i])
		value, _ := strconv.Atoi(values[i])
		result := models.Result{
			Value: uint(value),
			CriteriaID: uint(criteria),
			ExerciseID: item.ID,
		}
		db.DB.Create(&result)
	}

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