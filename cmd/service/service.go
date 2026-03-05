package service

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "service",
	Short: "Manage and test service connections",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Command.AddCommand(pingCmd)
	Command.AddCommand(listCmd)
}
