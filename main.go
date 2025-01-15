package main

import (
	"github.com/gofiber/fiber/v2"
)


func main() {
	app := fiber.New()

	app.Static("/", "./web/dist")

	app.Listen(":8080")
}
