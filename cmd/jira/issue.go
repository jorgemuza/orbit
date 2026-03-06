package jira

import "github.com/spf13/cobra"

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage Jira issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	issueCmd.AddCommand(issueListCmd)
	issueCmd.AddCommand(issueViewCmd)
	issueCmd.AddCommand(issueCreateCmd)
	issueCmd.AddCommand(issueEditCmd)
	issueCmd.AddCommand(issueAssignCmd)
	issueCmd.AddCommand(issueMoveCmd)
	issueCmd.AddCommand(issueDeleteCmd)
	issueCmd.AddCommand(issueCommentCmd)
	issueCmd.AddCommand(issueLinkCmd)
	issueCmd.AddCommand(issueUnlinkCmd)
	issueCmd.AddCommand(issueWorklogCmd)
	issueCmd.AddCommand(issueCloneCmd)
}
