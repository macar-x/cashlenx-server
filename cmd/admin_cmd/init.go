package admin_cmd

import (
	"fmt"

	"github.com/macar-x/cashlenx-server/service/manage_service"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize database with demo data",
	Long: `Initialize the database with demo categories and sample transactions.
Useful for testing and development.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := manage_service.InitializeDemoData()
		if err != nil {
			return err
		}

		fmt.Println("Database initialized with demo data successfully")
		fmt.Println("Demo data includes:")
		fmt.Println("  - 8 default categories")
		fmt.Println("  - 15 sample transactions")
		return nil
	},
}
