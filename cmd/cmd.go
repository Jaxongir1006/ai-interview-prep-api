package main

import (
	"os"

	"go-enterprise-blueprint/internal/app"

	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/spf13/cobra"
)

func main() {
	var root = &cobra.Command{}

	root.AddCommand(run())

	root.AddCommand(app.AuthCommands())
	root.AddCommand(app.MigrateCommands())
	// Add new modules CLI commands here...

	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// run registers a main command that runs all services.
func run() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run all services in one process",
		Run: func(_ *cobra.Command, _ []string) {
			err := app.Run()
			if err != nil {
				logger.Fatalx(err)
			}
		},
	}
}
