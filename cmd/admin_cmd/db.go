package admin_cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var adminToken string

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Database operations (admin only)",
	Long: `Database management operations requiring admin privileges.

Available sub-commands:
  connect  - Test database connection
  dump     - Create database dump
  restore  - Restore database from dump
  seed     - Seed database with demo data
  truncate - Clear all data from database`,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// adminToken is a shared flag for dangerous operations
		// It will be verified in each dangerous subcommand
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("must provide a valid sub command")
	},
}

func init() {
	// Add global admin-token flag for dangerous operations
	dbCmd.PersistentFlags().StringVarP(
		&adminToken, "admin-token", "t", "", "Admin token for dangerous operations")

	// Add all db subcommands
	dbCmd.AddCommand(connectCmd)
	dbCmd.AddCommand(dumpCmd)
	dbCmd.AddCommand(restoreDbCmd)
	dbCmd.AddCommand(seedCmd)
	dbCmd.AddCommand(truncateCmd)
}
