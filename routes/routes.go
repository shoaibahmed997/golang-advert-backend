package routes

import (
	"go-ecom/handler"
	"go-ecom/helper"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/login", handler.Login)
	api.Post("/signup", handler.Signup)
	api.Post("/change-password", helper.ParseToken, handler.ChangePassword)
	api.Get("/delete-account", helper.ParseToken, handler.DeleteUser)
}

func PostRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/posts", handler.GetAllPost)
	api.Get("/search/:searchterm", handler.SearchPost)
	api.Get("/posts/user/:email", handler.GetPostByUser)
	api.Get("/posts/category/:category", handler.GetPostByCategory)
	api.Get("/deletepost/:id", helper.ParseToken, handler.DeletePost)
	api.Get("/get20", handler.GetFirst20post)
	api.Post("/createpost", helper.ParseToken, handler.CreatePost)
}
