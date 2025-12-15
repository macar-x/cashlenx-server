package manage_cmd

import (
	"github.com/macar-x/cashlenx-server/service/manage_service"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "export data to excel",
	RunE: func(cmd *cobra.Command, args []string) error {
		return manage_service.ExportService(fromDate, toDate, filePath)
	},
}

func init() {
	exportCmd.Flags().StringVarP(&fromDate, "from", "f", "", "from date(include), e.x. 19700101")
	exportCmd.Flags().StringVarP(&toDate, "to", "t", "", "to date(include), e.x. 19700101")
	exportCmd.Flags().StringVarP(&filePath, "output", "o", "", "output path, default ./export.xlsx")
	ManageCmd.AddCommand(exportCmd)
}
