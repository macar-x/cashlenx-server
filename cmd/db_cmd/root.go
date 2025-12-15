package db_cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var DbCmd = &cobra.Command{
	Use:   "db",
	Short: "database operations",
	Long: `Database management operations.

Available sub-commands:
  connect - Test database connection
  migrate - Run database migrations
  seed    - Seed database with demo data`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("must provide a valid sub command")
	},
}
