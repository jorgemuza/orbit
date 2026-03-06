package github

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	ghsvc "github.com/jorgemuza/orbit/internal/service/github"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:     "run [subcommand]",
	Short:   "Manage GitHub Actions workflow runs",
	Aliases: []string{"actions"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var runListCmd = &cobra.Command{
	Use:   "list [owner/repo]",
	Short: "List workflow runs",
	Args:  cobra.ExactArgs(1),
	Example: `  orbit github run list octocat/hello-world
  orbit gh run list octocat/hello-world --branch main --status success`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		branch, _ := cmd.Flags().GetString("branch")
		status, _ := cmd.Flags().GetString("status")
		limit, _ := cmd.Flags().GetInt("limit")

		runs, err := client.ListWorkflowRuns(owner, repo, branch, status, limit)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(runs, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("%-12s %-12s %-12s %-20s %-10s %s\n", "ID", "STATUS", "CONCLUSION", "BRANCH", "EVENT", "CREATED")
		fmt.Printf("%-12s %-12s %-12s %-20s %-10s %s\n", "--", "------", "----------", "------", "-----", "-------")
		for _, r := range runs {
			created := ""
			if len(r.CreatedAt) >= 10 {
				created = r.CreatedAt[:10]
			}
			branch := r.HeadBranch
			if len(branch) > 18 {
				branch = branch[:15] + "..."
			}
			fmt.Printf("%-12d %-12s %-12s %-20s %-10s %s\n", r.ID, r.Status, r.Conclusion, branch, r.Event, created)
		}
		return nil
	},
}

var runViewCmd = &cobra.Command{
	Use:   "view [owner/repo] [run-id]",
	Short: "View a workflow run",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github run view octocat/hello-world 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid run ID: %s", args[1])
		}

		r, err := client.GetWorkflowRun(owner, repo, id)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("output")
		if format == "json" {
			data, _ := json.MarshalIndent(r, "", "  ")
			fmt.Println(string(data))
			return nil
		}

		fmt.Printf("ID:         %d\n", r.ID)
		fmt.Printf("Name:       %s\n", r.Name)
		fmt.Printf("Status:     %s\n", r.Status)
		fmt.Printf("Conclusion: %s\n", r.Conclusion)
		fmt.Printf("Branch:     %s\n", r.HeadBranch)
		fmt.Printf("SHA:        %s\n", r.HeadSHA)
		fmt.Printf("Event:      %s\n", r.Event)
		fmt.Printf("Created:    %s\n", r.CreatedAt)
		fmt.Printf("URL:        %s\n", r.HTMLURL)
		return nil
	},
}

var runCancelCmd = &cobra.Command{
	Use:   "cancel [owner/repo] [run-id]",
	Short: "Cancel a workflow run",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github run cancel octocat/hello-world 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid run ID: %s", args[1])
		}

		if err := client.CancelWorkflowRun(owner, repo, id); err != nil {
			return err
		}

		fmt.Printf("Canceled workflow run %d\n", id)
		return nil
	},
}

var runRerunCmd = &cobra.Command{
	Use:   "rerun [owner/repo] [run-id]",
	Short: "Re-run a workflow run",
	Args:  cobra.ExactArgs(2),
	Example: `  orbit github run rerun octocat/hello-world 12345`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("invalid run ID: %s", args[1])
		}

		if err := client.RerunWorkflowRun(owner, repo, id); err != nil {
			return err
		}

		fmt.Printf("Re-running workflow run %d\n", id)
		return nil
	},
}

var runWatchCmd = &cobra.Command{
	Use:   "watch [owner/repo] [run-id]",
	Short: "Watch a workflow run until it completes",
	Long: `Watch a workflow run, polling for status updates and displaying job progress.
If no run-id is provided, watches the most recent in-progress run.`,
	Args: cobra.RangeArgs(1, 2),
	Example: `  orbit github run watch jorgemuza/orbit
  orbit github run watch jorgemuza/orbit 12345
  orbit gh run watch jorgemuza/orbit --interval 10`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := resolveGitHubClient(cmd)
		if err != nil {
			return err
		}

		owner, repo, err := ghsvc.OwnerRepo(args[0])
		if err != nil {
			return err
		}

		interval, _ := cmd.Flags().GetInt("interval")

		var runID int
		if len(args) == 2 {
			runID, err = strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid run ID: %s", args[1])
			}
		} else {
			runs, err := client.ListWorkflowRuns(owner, repo, "", "in_progress", 1)
			if err != nil {
				return err
			}
			if len(runs) == 0 {
				// Try queued
				runs, err = client.ListWorkflowRuns(owner, repo, "", "queued", 1)
				if err != nil {
					return err
				}
			}
			if len(runs) == 0 {
				// Fall back to most recent
				runs, err = client.ListWorkflowRuns(owner, repo, "", "", 1)
				if err != nil {
					return err
				}
			}
			if len(runs) == 0 {
				return fmt.Errorf("no workflow runs found for %s/%s", owner, repo)
			}
			runID = runs[0].ID
		}

		r, err := client.GetWorkflowRun(owner, repo, runID)
		if err != nil {
			return err
		}

		fmt.Printf("Watching %s run #%d (%s) on %s...\n", r.Name, r.ID, r.Event, r.HeadBranch)
		fmt.Printf("URL: %s\n\n", r.HTMLURL)

		for {
			r, err = client.GetWorkflowRun(owner, repo, runID)
			if err != nil {
				return err
			}

			jobs, err := client.ListWorkflowRunJobs(owner, repo, runID, 100)
			if err != nil {
				return err
			}

			// Clear and reprint
			fmt.Printf("\033[2K\rStatus: %s", r.Status)
			if r.Conclusion != "" {
				fmt.Printf(" (%s)", r.Conclusion)
			}
			fmt.Println()

			for _, j := range jobs {
				icon := statusIcon(j.Status, j.Conclusion)
				elapsed := ""
				if j.StartedAt != "" {
					if start, err := time.Parse(time.RFC3339, j.StartedAt); err == nil {
						if j.CompletedAt != "" {
							if end, err := time.Parse(time.RFC3339, j.CompletedAt); err == nil {
								elapsed = formatDuration(end.Sub(start))
							}
						} else {
							elapsed = formatDuration(time.Since(start))
						}
					}
				}
				fmt.Printf("  %s %s", icon, j.Name)
				if elapsed != "" {
					fmt.Printf(" (%s)", elapsed)
				}
				fmt.Println()

				for _, s := range j.Steps {
					sIcon := statusIcon(s.Status, s.Conclusion)
					fmt.Printf("    %s %s\n", sIcon, s.Name)
				}
			}

			if r.Status == "completed" {
				fmt.Printf("\nRun %d finished: %s\n", r.ID, r.Conclusion)
				if r.Conclusion == "failure" || r.Conclusion == "cancelled" || r.Conclusion == "timed_out" {
					return fmt.Errorf("run concluded with: %s", r.Conclusion)
				}
				return nil
			}

			time.Sleep(time.Duration(interval) * time.Second)

			// Move cursor up to overwrite (status line + jobs + steps)
			lines := 1 // status line
			for _, j := range jobs {
				lines++ // job line
				lines += len(j.Steps)
			}
			for i := 0; i < lines; i++ {
				fmt.Print("\033[A\033[2K")
			}
		}
	},
}

func statusIcon(status, conclusion string) string {
	switch status {
	case "completed":
		switch conclusion {
		case "success":
			return "v"
		case "failure":
			return "X"
		case "cancelled":
			return "-"
		case "skipped":
			return "o"
		default:
			return "?"
		}
	case "in_progress":
		return "*"
	case "queued", "waiting", "pending":
		return "."
	default:
		return " "
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	}
	return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
}

func init() {
	runCmd.AddCommand(runListCmd)
	runCmd.AddCommand(runViewCmd)
	runCmd.AddCommand(runWatchCmd)
	runCmd.AddCommand(runCancelCmd)
	runCmd.AddCommand(runRerunCmd)

	runListCmd.Flags().String("branch", "", "filter by branch")
	runListCmd.Flags().String("status", "", "filter by status: completed, in_progress, queued")
	runListCmd.Flags().Int("limit", 20, "max results")

	runWatchCmd.Flags().Int("interval", 5, "polling interval in seconds")
}
