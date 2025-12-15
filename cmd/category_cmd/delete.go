package category_cmd

import (
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete category data",
	RunE: func(cmd *cobra.Command, args []string) error {
		return category_service.DeleteService(plainId, categoryName)
	},
}

func init() {
	deleteCmd.Flags().StringVarP(
		&plainId, "id", "i", "", "delete by id")
	deleteCmd.Flags().StringVarP(
		&categoryName, "name", "n", "", "delete by name")
	CategoryCmd.AddCommand(deleteCmd)
}
