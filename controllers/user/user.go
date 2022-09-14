package user

import (
	"phuhung273/progress-tracking/db"
	"phuhung273/progress-tracking/middleware"
	"phuhung273/progress-tracking/models"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func signUpView(c *fiber.Ctx) error {
	return c.Render("signup.html", fiber.Map{
		"title": "Sign Up",
	})
}

func signUp(c *fiber.Ctx) error {
	username := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	exist := db.DB.First(&user, "username = ?", username)
	if exist.RowsAffected > 0 {
		return c.Render("signup.html", fiber.Map{
			"title": "Sign Up",
			"error": "Username already existed",
		})
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	db.DB.Create(&models.User{ Username: username, Password: string(hash) })
	return c.Redirect("login")
}

func loginView(c *fiber.Ctx) error {
	return c.Render("login.html", fiber.Map{
		"title": "Login",
	})
}

func login(c *fiber.Ctx) error {
	username := c.FormValue("email")
	password := c.FormValue("password")

	var user models.User
	exist := db.DB.First(&user, "username = ?", username)
	if exist.RowsAffected == 0 {
		return c.Render("login.html", fiber.Map{
			"title": "Login",
			"error": "Username not existed",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return c.Render("login.html", fiber.Map{
			"title": "Login",
			"error": "Invalid username or password",
		})
	}

	sess, _ := middleware.SessionStore.Get(c)
	sess.Set("user_id", user.ID)
	sess.Save()
	
	return c.Redirect("/dashboard")
}

func signout(c *fiber.Ctx) error {
	sess, _ := middleware.SessionStore.Get(c)
	sess.Delete("user_id")
	sess.Save()
	
	return c.Redirect("login")
}

func Routing(router *fiber.App) {
	router.Route("/auth", func(router fiber.Router) {
		router.Get("/signup", signUpView)
		router.Post("/signup", signUp)
		router.Get("/login", loginView)
		router.Post("/login", login)
		router.Get("/logout", middleware.Auth, signout)
	})
}