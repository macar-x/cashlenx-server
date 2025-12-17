package db_cmd

import (
	"fmt"
	"time"

	"github.com/macar-x/cashlenx-server/service/manage_service"
	"github.com/macar-x/cashlenx-server/util"
	"github.com/spf13/cobra"
)

var dumpPath string

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "create database dump",
	Long: `Create a dump of all database data.
If no path is specified, creates dump in current directory with timestamp.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Verify ADMIN_TOKEN for dangerous operation
		if err := util.VerifyAdminToken(adminToken); err != nil {
			return err
		}

		if dumpPath == "" {
			dumpPath = fmt.Sprintf("cashlenx_dump_%s.json", time.Now().Format("20060102_150405"))
		}

		stats, err := manage_service.CreateBackup(dumpPath)
		if err != nil {
			return err
		}

		fmt.Printf("Database dump created successfully: %s\n", dumpPath)
		fmt.Println("\nStatistics:")
		fmt.Printf("  Categories: %d success, %d failed\n", stats.Categories.Success, stats.Categories.Failed)
		fmt.Printf("  Cash Flows: %d success, %d failed\n", stats.CashFlows.Success, stats.CashFlows.Failed)
		return nil
	},
}

func init() {
	dumpCmd.Flags().StringVarP(
		&dumpPath, "output", "o", "", "dump file path (optional, default: cashlenx_dump_TIMESTAMP.json)")

	DbCmd.AddCommand(dumpCmd)
}
