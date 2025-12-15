package category_cmd

import (
	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all categories",
	Long:  `List all categories in the system.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return category_service.ListAllService()
	},
}

func init() {
	CategoryCmd.AddCommand(listCmd)
}
