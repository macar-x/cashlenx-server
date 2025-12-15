package category_cmd

import (
	"fmt"

	"github.com/macar-x/cashlenx-server/service/category_service"
	"github.com/spf13/cobra"
)

var (
	categoryLimit  int
	categoryOffset int
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all categories",
	Long:  `List all categories in the system with optional pagination.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		categoryEntityList, _, err := category_service.ListAllService(categoryLimit, categoryOffset)
		if err != nil {
			return err
		}

		for index, categoryEntity := range categoryEntityList {
			fmt.Printf("category %d: %s (ID: %s)\n", index+categoryOffset, categoryEntity.Name, categoryEntity.Id.Hex())
		}

		return nil
	},
}

func init() {
	listCmd.Flags().IntVarP(
		&categoryLimit, "limit", "l", 50, "maximum number of records to return")
	listCmd.Flags().IntVarP(
		&categoryOffset, "offset", "o", 0, "number of records to skip")
	CategoryCmd.AddCommand(listCmd)
}
