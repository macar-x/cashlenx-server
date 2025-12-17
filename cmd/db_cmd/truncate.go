package db_cmd

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/macar-x/cashlenx-server/service/manage_service"
	"github.com/macar-x/cashlenx-server/util"
	"github.com/spf13/cobra"
)

var forceTruncate bool

var truncateCmd = &cobra.Command{
	Use:   "truncate",
	Short: "clear all data from database",
	Long: `Truncate (clear) all data from the database.
WARNING: This will delete all cash flows and categories permanently!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Verify ADMIN_TOKEN for dangerous operation
		if err := util.VerifyAdminToken(adminToken); err != nil {
			return err
		}

		if !forceTruncate {
			fmt.Println("WARNING: This will permanently delete ALL data from the database!")
			fmt.Print("Are you sure you want to continue? (yes/no): ")

			reader := bufio.NewReader(cmd.InOrStdin())
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response != "yes" && response != "y" {
				fmt.Println("Truncate cancelled")
				return nil
			}
		}

		stats, err := manage_service.TruncateDatabase()
		if err != nil {
			return err
		}

		fmt.Println("Database truncated successfully - all data has been cleared")
		fmt.Println("\nStatistics:")
		fmt.Printf("  Categories: %d success, %d failed\n", stats.Categories.Success, stats.Categories.Failed)
		fmt.Printf("  Cash Flows: %d success, %d failed\n", stats.CashFlows.Success, stats.CashFlows.Failed)
		return nil
	},
}

func init() {
	truncateCmd.Flags().BoolVarP(
		&forceTruncate, "force", "f", false, "skip confirmation prompt")

	DbCmd.AddCommand(truncateCmd)
}
