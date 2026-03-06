package jira

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var fieldListCmd = &cobra.Command{
	Use:   "field-list",
	Short: "List Jira fields (useful for finding custom field IDs)",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		var fields []map[string]any
		if err := client.DoGet("/rest/api/2/field", &fields); err != nil {
			return fmt.Errorf("listing fields: %w", err)
		}

		filter, _ := cmd.Flags().GetString("filter")
		filter = strings.ToLower(filter)

		for _, f := range fields {
			name := fmt.Sprintf("%v", f["name"])
			id := fmt.Sprintf("%v", f["id"])
			custom, _ := f["custom"].(bool)
			if filter != "" && !strings.Contains(strings.ToLower(name), filter) && !strings.Contains(strings.ToLower(id), filter) {
				continue
			}
			fmt.Printf("%-30s %-40s custom=%v\n", id, name, custom)
		}
		return nil
	},
}

func init() {
	fieldListCmd.Flags().String("filter", "", "filter fields by name (case-insensitive)")
}
