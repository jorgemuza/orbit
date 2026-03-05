package profile

import (
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "profile",
	Short: "Manage profiles",
	Long:  "Create, list, update, delete, and inspect profiles that group service connections.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	Command.AddCommand(createCmd)
	Command.AddCommand(listCmd)
	Command.AddCommand(showCmd)
	Command.AddCommand(deleteCmd)
	Command.AddCommand(setDefaultCmd)
	Command.AddCommand(addServiceCmd)
	Command.AddCommand(removeServiceCmd)
}
