package main

import (
	"go-ecom/database"
	"go-ecom/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	app := fiber.New()
	database.DBInit()
	app.Use(cors.New())
	routes.UserRoutes(app)
	routes.PostRoutes(app)
	app.Static("/", "./frontend")
	app.Static("/", "./public")
	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("./frontend/index.html")
	})

	defer database.DB.Close()

	err := app.Listen(":5000")
	if err != nil {
		log.Fatal("error", err)
	}
}
