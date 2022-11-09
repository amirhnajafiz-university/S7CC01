package cmd

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
)

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "api",
		Long: "command for starting the api service",
		Run: func(_ *cobra.Command, args []string) {
			main()
		},
	}
}

// main method of api service.
func main() {
	// creating a new fiber app
	app := fiber.New()

	// starting fiber
	if err := app.Listen(":5050"); err != nil {
		panic(err)
	}
}
