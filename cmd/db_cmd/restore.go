package db_cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/macar-x/cashlenx-server/service/manage_service"
	"github.com/spf13/cobra"
)

var (
	dbRestorePath  string
	forceDbRestore bool
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "restore database from dump",
	Long: `Restore database from a dump file.
WARNING: This will replace all existing data!`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dbRestorePath == "" {
			return errors.New("dump file path is required")
		}

		if !forceDbRestore {
			fmt.Println("WARNING: This will replace all existing data!")
			fmt.Print("Are you sure you want to continue? (yes/no): ")

			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response != "yes" && response != "y" {
				fmt.Println("Restore cancelled")
				return nil
			}
		}

		err := manage_service.RestoreBackup(dbRestorePath)
		if err != nil {
			return err
		}

		fmt.Printf("Database restored successfully from: %s\n", dbRestorePath)
		return nil
	},
}

func init() {
	restoreCmd.Flags().StringVarP(
		&dbRestorePath, "input", "i", "", "dump file path (required)")
	restoreCmd.Flags().BoolVarP(
		&forceDbRestore, "force", "f", false, "skip confirmation prompt")

	restoreCmd.MarkFlagRequired("input")
	DbCmd.AddCommand(restoreCmd)
}
