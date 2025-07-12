package router

import (
	"imageboard/controllers"

	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	main := router.Group("/")
	main.Get("/", controllers.HomePageController)
	main.Get("/login", controllers.LoginPageController)
	main.Get("/register", controllers.RegisterPageController)
	main.Get("/preferences", controllers.PreferencesPageController)

	posts := router.Group("/posts")
	posts.Get("/", controllers.PostsController)

	router.Use(controllers.NotFoundController)
}
