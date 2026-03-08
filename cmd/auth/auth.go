package auth

import (
	"fmt"

	"github.com/jorgemuza/orbit/cmd/cmdutil"
	"github.com/jorgemuza/orbit/internal/secrets"
	"github.com/spf13/cobra"
)

// Command is the top-level auth command.
var Command = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with 1Password (resolve and cache secrets)",
	Long: `Resolve all 1Password secret references (op://) from your config and cache
them locally. This triggers a single biometric prompt and caches the resolved
secrets for 8 hours, so subsequent orbit commands (including parallel ones)
work without additional prompts.

Run "orbit auth clear" to remove the cached secrets.`,
	Example: `  orbit auth
  orbit auth clear`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := cmdutil.LoadConfig(cmd)
		if err != nil {
			return err
		}

		// Collect all unique op:// references from all profiles
		refs := make(map[string]bool)
		for _, p := range cfg.Profiles {
			for _, svc := range p.Services {
				for _, val := range []string{svc.Auth.Token, svc.Auth.Username, svc.Auth.Password, svc.Auth.ClientSecret} {
					if secrets.IsSecretReference(val) {
						refs[val] = true
					}
				}
			}
		}

		if len(refs) == 0 {
			fmt.Println("No 1Password secret references found in config.")
			return nil
		}

		// Collect into a slice for batch resolution
		refList := make([]string, 0, len(refs))
		for ref := range refs {
			refList = append(refList, ref)
		}

		fmt.Printf("Resolving %d secret(s) from 1Password...\n", len(refList))

		resolved, err := secrets.ResolveAll(refList...)
		if err != nil {
			return fmt.Errorf("resolving secrets: %w", err)
		}

		// Build map for cache
		cacheMap := make(map[string]string, len(refList))
		for i, ref := range refList {
			cacheMap[ref] = resolved[i]
		}

		if err := secrets.SaveCache(cacheMap); err != nil {
			return fmt.Errorf("saving cache: %w", err)
		}

		fmt.Printf("Cached %d secret(s). Valid for %s.\n", len(cacheMap), secrets.CacheTTL)
		return nil
	},
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear cached secrets",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := secrets.ClearCache(); err != nil {
			return fmt.Errorf("clearing cache: %w", err)
		}
		fmt.Println("Secret cache cleared.")
		return nil
	},
}

func init() {
	Command.AddCommand(clearCmd)
}
