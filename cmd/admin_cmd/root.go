package admin_cmd

import (
	"github.com/spf13/cobra"
)

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Admin-only commands (requires admin privileges)",
	Long: `Admin-only commands that require admin authentication.
These commands are restricted and mirror the /api/admin/* endpoints.

Available sub-commands:
  db      - Database operations (connect, dump, restore, seed, truncate)
  manage  - Data management (backup, restore, import, export, stats, init, reset)`,
}

func init() {
	AdminCmd.AddCommand(dbCmd)
	AdminCmd.AddCommand(manageCmd)
}
