package db_cmd

import (
	"fmt"

	"github.com/macar-x/cashlenx-server/service/manage_service"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "seed database with demo data",
	Long: `Seed the database with demo categories and sample transactions.
This is an alias for 'manage init' command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := manage_service.InitializeDemoData()
		if err != nil {
			return err
		}

		fmt.Println("✅ Database seeded with demo data successfully")
		return nil
	},
}

func init() {
	DbCmd.AddCommand(seedCmd)
}
