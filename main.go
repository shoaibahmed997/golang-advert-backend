package main

import (
	"go-ecom/database"
	"go-ecom/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	app := fiber.New()
	database.DBInit()
	routes.UserRoutes(app)
	routes.PostRoutes(app)
	app.Static("/", "./public")

	defer database.DB.Close()

	err := app.Listen(":5000")
	if err != nil {
		log.Fatal("error", err)
	}
}
