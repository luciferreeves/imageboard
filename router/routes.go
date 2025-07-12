package router

import (
	"imageboard/controllers"

	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	main := router.Group("/")
	main.Get("/", controllers.HomeController)

	posts := router.Group("/posts")
	posts.Get("/", controllers.PostsController)

	// router.Get("/posts", controllers.PostsController)
	// router.Get("/register", controllers.RegisterController)
	// router.Get("/login", controllers.LoginController)
	// router.Get("/preferences", controllers.PreferencesController)

	router.Use(controllers.NotFoundController)
}
