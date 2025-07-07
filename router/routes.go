package router

import (
	"imageboard/controllers"

	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	router.Get("/", controllers.PostsController)
	router.Get("/register", controllers.RegisterController)
	router.Get("/login", controllers.LoginController)
	router.Get("/preferences", controllers.PreferencesController)

	router.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   "Not Found",
			"message": "The requested resource could not be found.",
			"status":  fiber.StatusNotFound,
		})
	})
}
