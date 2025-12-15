package category_cmd

import (
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query for category data",
	RunE: func(cmd *cobra.Command, args []string) error {
		return category_service.QueryService(plainId, parentPlainId, categoryName)
	},
}

func init() {
	queryCmd.Flags().StringVarP(
		&plainId, "id", "i", "", "query by id")
	queryCmd.Flags().StringVarP(
		&parentPlainId, "parent", "p", "", "query by parent id")
	queryCmd.Flags().StringVarP(
		&categoryName, "name", "n", "", "query by name")
	CategoryCmd.AddCommand(queryCmd)
}
