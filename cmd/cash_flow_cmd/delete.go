package cash_flow_cmd

import (
	"errors"
	"fmt"

	"github.com/macar-x/cashlenx-server/service/cash_flow_service"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete cash_flow by specific type",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Valid params through command.
		if cash_flow_service.IsDeleteFieldsConflicted(plainId, belongsDate) {
			return errors.New("should have one and only one delete type")
		}

		if plainId != "" {

			cashFlowEntity, err := cash_flow_service.DeleteById(plainId)
			if err != nil {
				return err
			}
			fmt.Println("cash_flow ", 0, ": ", cashFlowEntity.ToString())
			return nil
		}

		if belongsDate != "" {
			cashFlowEntityList, err := cash_flow_service.DeleteByDate(belongsDate)
			if err != nil {
				return err
			}
			if len(cashFlowEntityList) == 0 {
				fmt.Println("the day's flow is empty")
				return nil
			}
			for index, cashFlowEntity := range cashFlowEntityList {
				fmt.Println("cash_flow ", index, ": ", cashFlowEntity.ToString())
			}
			return nil
		}

		return errors.New("not supported delete type")
	},
}

func init() {
	deleteCmd.Flags().StringVarP(
		&plainId, "id", "i", "", "delete by id")
	deleteCmd.Flags().StringVarP(
		&belongsDate, "date", "b", "", "delete by belongs-date")
	CashCmd.AddCommand(deleteCmd)
}
