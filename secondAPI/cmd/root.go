package cmd

import "github.com/gofiber/fiber/v2"

// Execute
// main method of application.
func Execute() {
	// creating a new fiber app
	app := fiber.New()

	// starting fiber
	if err := app.Listen(":5050"); err != nil {
		panic(err)
	}
}
