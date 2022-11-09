package cmd

import "github.com/spf13/cobra"

func GetCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "processor",
		Long: "starting the processor service",
		Run: func(_ *cobra.Command, _ []string) {
			main()
		},
	}
}

// main method of processor.
func main() {
	// todo processor logic
}
