package secrets

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

// IsSecretReference checks if a value is a 1Password secret reference (op://vault/item/field).
func IsSecretReference(value string) bool {
	return strings.HasPrefix(value, "op://")
}

// Resolve resolves a value that may be a 1Password secret reference.
// If the value starts with "op://", it calls `op read` to retrieve the secret.
// Otherwise, it returns the value as-is.
func Resolve(value string) (string, error) {
	if !IsSecretReference(value) {
		return value, nil
	}
	return readSecret(value)
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
