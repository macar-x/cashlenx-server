package category_cmd

import (
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create new category",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := category_service.CreateService(parentPlainId, categoryName)
		return err
	},
}

func init() {
	createCmd.Flags().StringVarP(
		&parentPlainId, "parent", "p", "", "category's parent's id (optional)")
	createCmd.Flags().StringVarP(
		&categoryName, "name", "n", "", "category's name (required)")
	CategoryCmd.AddCommand(createCmd)
}
