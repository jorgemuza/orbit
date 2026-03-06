package secrets

import "testing"

func TestIsSecretReference(t *testing.T) {
	tests := []struct {
		value string
		want  bool
	}{
		{"op://vault/item/field", true},
		{"op://Dev/jira-token/credential", true},
		{"plain-token", false},
		{"", false},
		{"OP://uppercase", false},
	}
	for _, tt := range tests {
		if got := IsSecretReference(tt.value); got != tt.want {
			t.Errorf("IsSecretReference(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}

func TestResolvePlainValue(t *testing.T) {
	val, err := Resolve("plain-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "plain-token" {
		t.Fatalf("expected 'plain-token', got %q", val)
	}
}

func TestResolveEmptyValue(t *testing.T) {
	val, err := Resolve("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Fatalf("expected empty string, got %q", val)
	}
}
