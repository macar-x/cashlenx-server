package admin_cmd

import (
	"github.com/macar-x/cashlenx-server/service/manage_service"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import data from Excel",
	Long:  `Import cash flow data from Excel file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return manage_service.ImportService(filePath)
	},
}

func init() {
	importCmd.Flags().StringVarP(&filePath, "input", "i", "", "input path, e.g. ~/export.xlsx")
}
