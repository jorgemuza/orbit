package service

import (
	"github.com/paybook/aidlc-cli/cmd/cmdutil"
	"github.com/paybook/aidlc-cli/internal/output"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List services in the active profile",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, p, err := cmdutil.ResolveProfile(cmd)
		if err != nil {
			return err
		}

		headers, rowFn := output.ServiceTable(p.Services)
		return output.Print(cmdutil.OutputFormat(cmd), p.Services, headers, rowFn)
	},
}
