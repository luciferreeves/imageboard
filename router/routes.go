package router

import (
	"imageboard/controllers"

	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	main := router.Group("/")
	main.Get("/", controllers.HomePageController)
	main.Get("/register", controllers.RegisterPageController)
	main.Get("/preferences", controllers.PreferencesPageController)

	posts := router.Group("/posts")
	posts.Get("/", controllers.PostsController)

	login := router.Group("/login")
	login.Get("/", controllers.LoginPageController)
	login.Post("/", controllers.LoginPostController)

	router.Use(controllers.NotFoundController)
}
