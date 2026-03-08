package secrets

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// CacheTTL is how long cached secrets remain valid.
const CacheTTL = 8 * time.Hour

// CacheEntry is the on-disk format for the secret cache.
type CacheEntry struct {
	ExpiresAt time.Time         `json:"expires_at"`
	Secrets   map[string]string `json:"secrets"`
}

// cachePath returns ~/.config/orbit/.secret-cache
func cachePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "orbit", ".secret-cache"), nil
}

// LoadCache reads the secret cache from disk. Returns nil if the cache
// doesn't exist, is expired, or can't be read.
func LoadCache() *CacheEntry {
	path, err := cachePath()
	if err != nil {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil
	}
	if time.Now().After(entry.ExpiresAt) {
		os.Remove(path)
		return nil
	}
	return &entry
}

// SaveCache merges resolved secrets into the existing cache file and writes
// it back with a refreshed TTL.
func SaveCache(secrets map[string]string) error {
	path, err := cachePath()
	if err != nil {
		return err
	}
	// Merge with existing cache so we don't discard previously cached secrets
	merged := make(map[string]string)
	if existing := LoadCache(); existing != nil {
		for k, v := range existing.Secrets {
			merged[k] = v
		}
	}
	for k, v := range secrets {
		merged[k] = v
	}
	entry := CacheEntry{
		ExpiresAt: time.Now().Add(CacheTTL),
		Secrets:   merged,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

// ClearCache removes the secret cache file.
func ClearCache() error {
	path, err := cachePath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

var fileCacheOnce sync.Once

// ensureFileCache lazily loads the on-disk cache into the in-memory cache
// on the first call to Resolve or ResolveAll, avoiding file I/O on commands
// that don't need secrets (e.g., orbit version, orbit --help).
func ensureFileCache() {
	fileCacheOnce.Do(func() {
		entry := LoadCache()
		if entry == nil {
			return
		}
		for ref, val := range entry.Secrets {
			cache.Store(ref, val)
		}
	})
}
