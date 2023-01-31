package cmd

import (
	api "github.com/ceit-aut/S7CC01/internal/api/cmd"
	processor "github.com/ceit-aut/S7CC01/internal/processor/cmd"
	"github.com/ceit-aut/S7CC01/internal/test"

	"github.com/spf13/cobra"
)

// Execute
// services with golang cobra.
func Execute() {
	cmd := cobra.Command{}

	cmd.AddCommand(
		api.GetCommand(),
		processor.GetCommand(),
		test.GetTestCommands(),
	)

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
