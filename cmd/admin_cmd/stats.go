package admin_cmd

import (
	"fmt"

	"github.com/macar-x/cashlenx-server/service/manage_service"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show database statistics",
	Long:  `Display statistics about the database including record counts and storage info.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stats, err := manage_service.GetDatabaseStats()
		if err != nil {
			return err
		}

		fmt.Println("\n=== Database Statistics ===")
		fmt.Printf("Cash Flow Records:  %d\n", stats.CashFlowCount)
		fmt.Printf("  - Income:         %d\n", stats.IncomeCount)
		fmt.Printf("  - Expense:        %d\n", stats.ExpenseCount)
		fmt.Printf("Categories:         %d\n", stats.CategoryCount)
		fmt.Printf("\nFinancial Summary:\n")
		fmt.Printf("  Total Income:     %.2f\n", stats.TotalIncome)
		fmt.Printf("  Total Expense:    %.2f\n", stats.TotalExpense)
		fmt.Printf("  Balance:          %.2f\n", stats.Balance)
		fmt.Printf("\nDate Range:\n")
		fmt.Printf("  Earliest:         %s\n", stats.EarliestDate)
		fmt.Printf("  Latest:           %s\n", stats.LatestDate)

		return nil
	},
}
