package controllers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func ListResponse(c *fiber.Ctx, items interface{}, count int64) error {
	c.Set("Access-Control-Expose-Headers", "X-Total-Count")
	c.Set("X-Total-Count", fmt.Sprint(count))
	return c.JSON(items)
}