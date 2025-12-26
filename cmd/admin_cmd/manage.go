package admin_cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	fromDate string
	toDate   string
	filePath string
)

var manageCmd = &cobra.Command{
	Use:   "manage",
	Short: "Data management and utilities (admin only)",
	Long: `Manage application data including import, export, backup, and restore.

Available sub-commands:
  export  - Export data to Excel
  import  - Import data from Excel
  backup  - Create database backup
  restore - Restore from backup
  init    - Initialize with demo data
  reset   - Clear all data (dangerous!)
  stats   - Show database statistics
  indexes - Manage database indexes`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("must provide a valid sub command")
	},
}

func init() {
	// Add all manage subcommands
	manageCmd.AddCommand(exportCmd)
	manageCmd.AddCommand(importCmd)
	manageCmd.AddCommand(backupCmd)
	manageCmd.AddCommand(restoreBackupCmd)
	manageCmd.AddCommand(initCmd)
	manageCmd.AddCommand(resetCmd)
	manageCmd.AddCommand(statsCmd)
	manageCmd.AddCommand(indexesCmd)
}
