package home

import (
	"database/sql"
	"phuhung273/progress-tracking/db"
	"phuhung273/progress-tracking/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Result struct {
	Category string
	Criteria string
	Value uint
	Date time.Time
}

func index(c *fiber.Ctx) error {
	sess, _ := middleware.SessionStore.Get(c)

	userId := sess.Get("user_id")

	labels := []string{}
	values := []uint{}

	results, _ := db.DB.Raw(`
		SELECT criterias.name AS criteria, categories.name AS category, date(exercises.created_at) AS date, SUM(value) AS value
		FROM results 
		INNER JOIN settings AS criterias ON results.criteria_id = criterias.id
		INNER JOIN exercises ON results.exercise_id = exercises.id
		INNER JOIN settings AS categories ON exercises.category_id = categories.id
		WHERE categories.name = @readingCategory AND exercises.user_id = @userId
		GROUP BY criterias.name, categories.name, date(exercises.created_at)
	`, 
		sql.Named("readingCategory", "Reading"), sql.Named("userId", userId)).Rows()
		
	for results.Next() {
		var result Result
		db.DB.ScanRows(results, &result)
		labels = append(labels, result.Date.Format("02-01"))
		values = append(values, result.Value)
	}

	return c.Render("home.html", fiber.Map{
		"title": "Home",
		"labels": labels,
		"values": values,
	})
}

func Routing(router *fiber.App) {
	router.Get("/dashboard", middleware.Auth, index)
}