package db_cmd

import (
	"fmt"

	"github.com/macar-x/cashlenx-server/service/db_service"
	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "test database connection",
	Long:  `Test the connection to the configured database and display connection info.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := db_service.TestConnection()
		if err != nil {
			fmt.Println("❌ Database connection failed")
			return err
		}

		fmt.Println("✅ Database connection successful")
		fmt.Printf("\nConnection Info:\n")
		fmt.Printf("  Type:     %s\n", info.Type)
		fmt.Printf("  Host:     %s\n", info.Host)
		fmt.Printf("  Database: %s\n", info.Database)
		fmt.Printf("  Status:   %s\n", info.Status)

		return nil
	},
}

func init() {
	DbCmd.AddCommand(connectCmd)
}
