package manage_cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var (
	fromDate string
	toDate   string
	filePath string
)

var ManageCmd = &cobra.Command{
	Use:   "manage",
	Short: "data management and utilities",
	Long: `Manage application data including import, export, backup, and restore.

Available sub-commands:
  export  - Export data to Excel
  import  - Import data from Excel
  backup  - Create database backup
  restore - Restore from backup
  init    - Initialize with demo data
  reset   - Clear all data (dangerous!)
  stats   - Show database statistics`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("must provide a valid sub command")
	},
}
