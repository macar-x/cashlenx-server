package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/macar-x/cashlenx-server/cmd/cash_flow_cmd"
	"github.com/macar-x/cashlenx-server/cmd/category_cmd"
	"github.com/macar-x/cashlenx-server/cmd/db_cmd"
	"github.com/macar-x/cashlenx-server/cmd/manage_cmd"
	"github.com/macar-x/cashlenx-server/cmd/server_cmd"
	"github.com/macar-x/cashlenx-server/util"
	"github.com/macar-x/cashlenx-server/util/database"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cashlenx",
	Short: "Personal finance management - See your money clearly",
	Long: `CashLenX - Personal finance management application
	
Track your daily cash flow, manage categories, and gain insights into your spending habits.
Use 'cashlenx [command] --help' for more information about a command.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize database connection pool
		dbType := util.GetConfigByKey("db.type")
		if dbType == "mongodb" {
			if err := database.InitMongoDbConnection(); err != nil {
				util.Logger.Errorw("Failed to initialize MongoDB connection", "error", err)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CashLenX - See your money clearly")
		fmt.Println("Use 'cashlenx --help' for available commands")
	},
}

func Execute() {
	// Setup graceful shutdown
	setupGracefulShutdown()

	cobra.CheckErr(rootCmd.Execute())
}

func setupGracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		util.Logger.Info("Shutdown signal received, cleaning up...")

		// Close database connections
		dbType := util.GetConfigByKey("db.type")
		if dbType == "mongodb" {
			database.ShutdownMongoDbConnection()
		}

		util.Logger.Info("Cleanup complete, exiting")
		os.Exit(0)
	}()
}

func init() {
	rootCmd.AddCommand(server_cmd.ServerCmd)
	rootCmd.AddCommand(cash_flow_cmd.CashCmd)
	rootCmd.AddCommand(category_cmd.CategoryCmd)
	rootCmd.AddCommand(manage_cmd.ManageCmd)
	rootCmd.AddCommand(db_cmd.DbCmd)
}
