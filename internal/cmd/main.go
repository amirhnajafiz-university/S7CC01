package cmd

import (
	api "github.com/ceit-aut/ad-registration-service/internal/api/cmd"
	processor "github.com/ceit-aut/ad-registration-service/internal/processor/cmd"
	"github.com/spf13/cobra"
)

// Execute
// services with golang cobra.
func Execute() {
	cmd := cobra.Command{}

	cmd.AddCommand(
		api.GetCommand(),
		processor.GetCommand(),
	)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
