package server_cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "manage api server",
	Long: `
Managing application API server.
Provide sub-commands: [start].`,

	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("must provide a valid sub command")
	},
}
