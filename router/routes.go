package router

import (
	"imageboard/controllers"

	"github.com/gofiber/fiber/v2"
)

func Initialize(router *fiber.App) {
	router.Static("/static", "./static")

	main := router.Group("/")
	main.Get("/", controllers.HomePageController)

	posts := router.Group("/posts")
	posts.Get("/", controllers.PostsPageController)
	posts.Get("/new", controllers.PostsUploadPageController)
	posts.Post("/new", controllers.PostsUploadPostController)
	posts.Get("/new/ilinkfetch", controllers.PostsUploadImageLinkProxyController)
	posts.Get("/:id", controllers.PostsSinglePostPageController)
	posts.Post("/:id/favourite", controllers.PostsSinglePostFavouritePostController)
	posts.Get("/:id/edit", controllers.PostsSinglePostEditPageController)
	posts.Post("/:id/edit", controllers.PostsSinglePostEditPostController)

	tags := router.Group("/tags")
	tags.Get("/search.json", controllers.TagsSearchJSONController)
	tags.Post("/create.json", controllers.FindOrCreateTagJSONController)
	tags.Get("/search_for_image.json", controllers.TagsSearchForImageJSONController)
	tags.Post("/add_to_image.json", controllers.TagsAddToImageJSONController)
	tags.Post("/remove_from_image.json", controllers.TagsRemoveFromImageJSONController)

	login := router.Group("/login")
	login.Get("/", controllers.LoginPageController)
	login.Post("/", controllers.LoginPostController)

	logout := router.Group("/logout")
	logout.Get("/", controllers.LogoutController)

	register := router.Group("/register")
	register.Get("/", controllers.RegisterPageController)
	register.Post("/", controllers.RegisterPostController)

	account := router.Group("/account")
	account.Get("/verify", controllers.VerifyEmailController)

	preferences := router.Group("/preferences")
	preferences.Get("/", controllers.PreferencesPageController)

	router.Use(controllers.NotFoundController)
}
