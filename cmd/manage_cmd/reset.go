package manage_cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/macar-x/cashlenx-server/service/manage_service"
	"github.com/spf13/cobra"
)

var forceReset bool

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "clear all database data",
	Long: `Clear all data from the database.
WARNING: This operation cannot be undone. Create a backup first!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !forceReset {
			fmt.Println("⚠️  WARNING: This will DELETE ALL DATA from the database!")
			fmt.Println("This operation cannot be undone.")
			fmt.Print("Type 'DELETE ALL' to confirm: ")

			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(response)

			if response != "DELETE ALL" {
				fmt.Println("Reset cancelled")
				return nil
			}
		}

		err := manage_service.ResetDatabase()
		if err != nil {
			return err
		}

		fmt.Println("✅ Database reset successfully - all data cleared")
		return nil
	},
}

func init() {
	resetCmd.Flags().BoolVarP(
		&forceReset, "force", "f", false, "skip confirmation prompt (dangerous!)")

	ManageCmd.AddCommand(resetCmd)
}
