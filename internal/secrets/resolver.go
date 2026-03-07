package secrets

import (
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

// cache stores resolved secrets for the lifetime of the process so that
// repeated references to the same op:// URI only trigger a single `op read`
// call (and therefore a single biometric prompt).
var cache sync.Map

// Resolve resolves a value that may be a 1Password secret reference.
// If the value starts with "op://", it calls `op read` to retrieve the secret.
// Resolved values are cached in-process to avoid repeated biometric prompts.
// Otherwise, it returns the value as-is.
func Resolve(value string) (string, error) {
	if !IsSecretReference(value) {
		return value, nil
	}

	// Check cache first
	if cached, ok := cache.Load(value); ok {
		return cached.(string), nil
	}

	secret, err := readSecret(value)
	if err != nil {
		return "", err
	}

	cache.Store(value, secret)
	return secret, nil
}

func readSecret(reference string) (string, error) {
	out, err := exec.Command("op", "read", reference).Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("1password: op read failed: %s", strings.TrimSpace(string(exitErr.Stderr)))
		}
		return "", fmt.Errorf("1password: op CLI not found or not configured: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}
