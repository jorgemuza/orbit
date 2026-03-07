package jira

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var screenCmd = &cobra.Command{
	Use:   "screen",
	Short: "Manage Jira screens (add/remove fields from create/edit/view screens)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var screenListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all screens",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		screens, err := client.ListScreens(200)
		if err != nil {
			return err
		}

		filter, _ := cmd.Flags().GetString("filter")
		filter = strings.ToLower(filter)

		for _, s := range screens {
			if filter != "" && !strings.Contains(strings.ToLower(s.Name), filter) {
				continue
			}
			desc := ""
			if s.Description != "" {
				desc = "  — " + s.Description
			}
			fmt.Printf("%-8d %s%s\n", s.ID, s.Name, desc)
		}
		return nil
	},
}

var screenTabListCmd = &cobra.Command{
	Use:   "tab-list [screen-id]",
	Short: "List tabs on a screen",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		screenID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid screen ID: %s", args[0])
		}

		tabs, err := client.ListScreenTabs(screenID)
		if err != nil {
			return err
		}

		for _, t := range tabs {
			fmt.Printf("%-8d %s\n", t.ID, t.Name)
		}
		return nil
	},
}

var screenFieldListCmd = &cobra.Command{
	Use:   "field-list [screen-id] [tab-id]",
	Short: "List fields on a screen tab",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		screenID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid screen ID: %s", args[0])
		}
		tabID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid tab ID: %s", args[1])
		}

		fields, err := client.ListScreenTabFields(screenID, tabID)
		if err != nil {
			return err
		}

		for _, f := range fields {
			fmt.Printf("%-30s %s\n", f.ID, f.Name)
		}
		return nil
	},
}

var screenFieldAddCmd = &cobra.Command{
	Use:   "field-add [screen-id] [tab-id]",
	Short: "Add fields to a screen tab",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		screenID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid screen ID: %s", args[0])
		}
		tabID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid tab ID: %s", args[1])
		}

		fieldIDs, _ := cmd.Flags().GetStringSlice("fields")
		if len(fieldIDs) == 0 {
			return fmt.Errorf("--fields is required (comma-separated field IDs)")
		}

		for _, fieldID := range fieldIDs {
			fieldID = strings.TrimSpace(fieldID)
			if err := client.AddFieldToScreen(screenID, tabID, fieldID); err != nil {
				fmt.Printf("Failed to add %s: %v\n", fieldID, err)
			} else {
				fmt.Printf("Added %s to screen %d tab %d\n", fieldID, screenID, tabID)
			}
		}
		return nil
	},
}

var screenTabCreateCmd = &cobra.Command{
	Use:   "tab-create [screen-id] [tab-name]",
	Short: "Create a new tab on a screen",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		screenID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid screen ID: %s", args[0])
		}

		tab, err := client.CreateScreenTab(screenID, args[1])
		if err != nil {
			return err
		}

		fmt.Printf("Created tab %d (%s) on screen %d\n", tab.ID, tab.Name, screenID)
		return nil
	},
}

var screenFieldRemoveCmd = &cobra.Command{
	Use:   "field-remove [screen-id] [tab-id]",
	Short: "Remove fields from a screen tab",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		screenID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid screen ID: %s", args[0])
		}
		tabID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid tab ID: %s", args[1])
		}

		fieldIDs, _ := cmd.Flags().GetStringSlice("fields")
		if len(fieldIDs) == 0 {
			return fmt.Errorf("--fields is required (comma-separated field IDs)")
		}

		for _, fieldID := range fieldIDs {
			fieldID = strings.TrimSpace(fieldID)
			if err := client.RemoveFieldFromScreen(screenID, tabID, fieldID); err != nil {
				fmt.Printf("Failed to remove %s: %v\n", fieldID, err)
			} else {
				fmt.Printf("Removed %s from screen %d tab %d\n", fieldID, screenID, tabID)
			}
		}
		return nil
	},
}

var screenFieldMoveCmd = &cobra.Command{
	Use:   "field-move [screen-id] [source-tab-id] [target-tab-id]",
	Short: "Move fields from one tab to another",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveJiraClient(cmd)
		if err != nil {
			return err
		}

		screenID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid screen ID: %s", args[0])
		}
		sourceTabID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid source tab ID: %s", args[1])
		}
		targetTabID, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("invalid target tab ID: %s", args[2])
		}

		fieldIDs, _ := cmd.Flags().GetStringSlice("fields")
		if len(fieldIDs) == 0 {
			return fmt.Errorf("--fields is required (comma-separated field IDs)")
		}

		for _, fieldID := range fieldIDs {
			fieldID = strings.TrimSpace(fieldID)
			if err := client.MoveFieldToTab(screenID, sourceTabID, targetTabID, fieldID); err != nil {
				fmt.Printf("Failed to move %s: %v\n", fieldID, err)
			} else {
				fmt.Printf("Moved %s from tab %d to tab %d on screen %d\n", fieldID, sourceTabID, targetTabID, screenID)
			}
		}
		return nil
	},
}

func init() {
	screenListCmd.Flags().String("filter", "", "filter screens by name (case-insensitive)")
	screenFieldAddCmd.Flags().StringSlice("fields", nil, "field IDs to add (comma-separated)")
	screenFieldRemoveCmd.Flags().StringSlice("fields", nil, "field IDs to remove (comma-separated)")
	screenFieldMoveCmd.Flags().StringSlice("fields", nil, "field IDs to move (comma-separated)")

	screenCmd.AddCommand(screenListCmd)
	screenCmd.AddCommand(screenTabListCmd)
	screenCmd.AddCommand(screenTabCreateCmd)
	screenCmd.AddCommand(screenFieldListCmd)
	screenCmd.AddCommand(screenFieldAddCmd)
	screenCmd.AddCommand(screenFieldRemoveCmd)
	screenCmd.AddCommand(screenFieldMoveCmd)
}
