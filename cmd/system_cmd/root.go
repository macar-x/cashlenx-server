package system_cmd

import (
	"github.com/spf13/cobra"
)

var SystemCmd = &cobra.Command{
	Use:   "system",
	Short: "System related commands",
	Long:  "System related commands like health check and version info",
}

func init() {
	SystemCmd.AddCommand(VersionCmd)
	SystemCmd.AddCommand(HealthCmd)
}
