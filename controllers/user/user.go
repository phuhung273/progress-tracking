package user

import (
	"os"
	"phuhung273/progress-tracking/db"
	"phuhung273/progress-tracking/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// func signUpView(c *fiber.Ctx) error {
// 	return c.Render("signup.html", fiber.Map{
// 		"title": "Sign Up",
// 	})
// }

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

// func loginView(c *fiber.Ctx) error {
// 	return c.Render("login.html", fiber.Map{
// 		"title": "Login",
// 	})
// }
type LoginRequest struct {
    Username string `json:"username" form:"username"`
    Password string `json:"password" form:"password"`
}

func login(c *fiber.Ctx) error {
	loginRequest := new(LoginRequest)

	if err := c.BodyParser(loginRequest); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"message": err,
		})
	}

	var user models.User
	exist := db.DB.First(&user, "username = ?", loginRequest.Username)
	if exist.RowsAffected == 0 {
		// return c.Render("login.html", fiber.Map{
		// 	"title": "Login",
		// 	"error": "Username not existed",
		// })
		return c.Status(401).JSON(fiber.Map{
			"message": "You are not registered",
		})
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		// return c.Render("login.html", fiber.Map{
		// 	"title": "Login",
		// 	"error": "Invalid username or password",
		// })
		return c.Status(401).JSON(fiber.Map{
			"message": "Invalid username or password",
		})
	}

	// sess, _ := middleware.SessionStore.Get(c)
	// sess.Set("user_id", user.ID)
	// sess.Save()

	exp, _ := strconv.Atoi(os.Getenv("JWT_EXP"))
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp": time.Now().Add(time.Hour * time.Duration(exp)).Unix(),
	})
	
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	
	return c.JSON(fiber.Map{
		"user_id": user.ID,
		"access_token": tokenString,
	})
}

// func signout(c *fiber.Ctx) error {
// 	// sess, _ := middleware.SessionStore.Get(c)
// 	// sess.Delete("user_id")
// 	// sess.Save()
	
// 	return c.Redirect("login")
// }

func Routing(router *fiber.App) {
	router.Route("/auth", func(router fiber.Router) {
		// router.Get("/signup", signUpView)
		router.Post("/signup", signUp)
		// router.Get("/login", loginView)
		router.Post("/login", login)
		// router.Get("/logout", middleware.Auth, signout)
	})
}