package app

import (
	"github.com/code19m/errx"
	"github.com/spf13/cobra"
)

func MigrateCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Database migration commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "up",
		Short: "Apply all pending database migrations",
		RunE: func(_ *cobra.Command, _ []string) error {
			app := newApp()
			defer app.shutdownInfraComponents()

			err := app.initSharedComponents()
			if err != nil {
				return errx.Wrap(err)
			}
			return app.migrateUp()
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "down",
		Short: "Rollback last applied database migration",
		RunE: func(_ *cobra.Command, _ []string) error {
			app := newApp()
			defer app.shutdownInfraComponents()

			err := app.initSharedComponents()
			if err != nil {
				return errx.Wrap(err)
			}
			return app.migrateDown()
		},
	})

	return cmd
}

func AuthCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Auth module CLI commands",
	}

	var username, password string

	createSuperadminCmd := &cobra.Command{
		Use:   "create-superadmin",
		Short: "Create superadmin account for system bootstrap",
		RunE: func(_ *cobra.Command, _ []string) error {
			app := newApp()
			defer app.shutdownInfraComponents()

			err := app.init()
			if err != nil {
				return errx.Wrap(err)
			}

			return app.auth.CreateSuperadmin(username, password)
		},
	}

	createSuperadminCmd.Flags().StringVar(&username, "username", "", "Admin username (skips interactive prompt)")
	createSuperadminCmd.Flags().StringVar(&password, "password", "", "Admin password (skips interactive prompt)")

	cmd.AddCommand(createSuperadminCmd)
	// Add auth modules new CLI commands here...

	return cmd
}

// Add your new CLI commands here...
