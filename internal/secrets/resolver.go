package secrets

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// IsSecretReference checks if a value is a 1Password secret reference (op://vault/item/field).
func IsSecretReference(value string) bool {
	return strings.HasPrefix(value, "op://")
}

// HasSecretReferences returns true if any of the given values is a secret reference.
func HasSecretReferences(values ...string) bool {
	for _, v := range values {
		if IsSecretReference(v) {
			return true
		}
	}
	return false
}

// cache stores resolved secrets for the lifetime of the process so that
// repeated references to the same op:// URI don't trigger additional prompts.
var cache sync.Map

// Resolve resolves a single value that may be a 1Password secret reference.
// Prefer ResolveAll for multiple references — it batches them into one `op`
// call, requiring only a single biometric prompt.
func Resolve(value string) (string, error) {
	resolved, err := ResolveAll(value)
	if err != nil {
		return "", err
	}
	return resolved[0], nil
}

// ResolveAll resolves multiple values that may contain 1Password secret
// references. All op:// references are batched into a single `op inject`
// call so that only one biometric prompt is needed regardless of how many
// secrets are being resolved.
func ResolveAll(values ...string) ([]string, error) {
	ensureFileCache()
	results := make([]string, len(values))

	// Identify which values need resolution and which are already cached
	type pending struct {
		index int
		ref   string
	}
	var toResolve []pending

	for i, v := range values {
		if !IsSecretReference(v) {
			results[i] = v
			continue
		}
		if cached, ok := cache.Load(v); ok {
			results[i] = cached.(string)
			continue
		}
		toResolve = append(toResolve, pending{index: i, ref: v})
	}

	if len(toResolve) == 0 {
		return results, nil
	}

	// Single reference — use op read directly (simpler)
	if len(toResolve) == 1 {
		secret, err := readSecret(toResolve[0].ref)
		if err != nil {
			return nil, err
		}
		cache.Store(toResolve[0].ref, secret)
		results[toResolve[0].index] = secret
		return results, nil
	}

	// Multiple references — batch via op inject (single biometric prompt).
	// We build a JSON template where keys map to op:// references, pipe it
	// through `op inject`, and parse the resolved JSON back.
	template := make(map[string]string, len(toResolve))
	for i, p := range toResolve {
		key := fmt.Sprintf("k%d", i)
		template[key] = "{{ " + p.ref + " }}"
	}

	tmplJSON, err := json.Marshal(template)
	if err != nil {
		return nil, fmt.Errorf("1password: building template: %w", err)
	}

	resolved, err := opInject(string(tmplJSON))
	if err != nil {
		return nil, err
	}

	var resolvedMap map[string]string
	if err := json.Unmarshal([]byte(resolved), &resolvedMap); err != nil {
		return nil, fmt.Errorf("1password: parsing inject output: %w", err)
	}

	for i, p := range toResolve {
		key := fmt.Sprintf("k%d", i)
		secret, ok := resolvedMap[key]
		if !ok {
			return nil, fmt.Errorf("1password: missing resolved value for %s", p.ref)
		}
		cache.Store(p.ref, secret)
		results[p.index] = secret
	}

	return results, nil
}

// opExecError wraps an exec.Command error with 1Password-specific context.
func opExecError(action string, err error) error {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return fmt.Errorf("1password: op %s failed: %s", action, strings.TrimSpace(string(exitErr.Stderr)))
	}
	return fmt.Errorf("1password: op CLI not found or not configured: %w", err)
}

func readSecret(reference string) (string, error) {
	out, err := exec.Command("op", "read", reference).Output()
	if err != nil {
		return "", opExecError("read", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func opInject(template string) (string, error) {
	cmd := exec.Command("op", "inject")
	cmd.Stdin = strings.NewReader(template)
	out, err := cmd.Output()
	if err != nil {
		return "", opExecError("inject", err)
	}
	return strings.TrimSpace(string(out)), nil
}
