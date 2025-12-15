package cash_flow_cmd

import (
	"errors"
	"fmt"

	"github.com/macar-x/cashlenx-server/service/cash_flow_service"
	"github.com/spf13/cobra"
)

var outcomeCmd = &cobra.Command{
	Use:   "outcome",
	Short: "add new outcome cash_flow",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cash_flow_service.IsOutcomeRequiredFiledSatisfied(categoryName, amount) {
			return errors.New("some required fields are empty")
		}
		cashFlowEntity, err := cash_flow_service.SaveOutcome(belongsDate, categoryName, amount, descriptionExact)
		if err != nil {
			return err
		}
		fmt.Println("cash_flow ", 0, ": ", cashFlowEntity.ToString())
		return nil
	},
}

func init() {
	outcomeCmd.Flags().StringVarP(
		&belongsDate, "date", "b", "", "flow's belongs-date (optional, blank for today)")
	outcomeCmd.Flags().StringVarP(
		&categoryName, "category", "c", "", "flow's category name (required)")
	outcomeCmd.Flags().Float64VarP(
		&amount, "amount", "a", 0.00, "flow's amount (required)")
	outcomeCmd.Flags().StringVarP(
		&descriptionExact, "description", "d", "", "flow's description (optional, could be blank)")
	CashCmd.AddCommand(outcomeCmd)
}
