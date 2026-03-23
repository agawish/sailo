package creds

import "os"

// FilterEnv returns only the environment variables that are in the allowlist.
// This prevents accidental leakage of secrets into workspace containers.
func FilterEnv(allowlist []string) map[string]string {
	result := make(map[string]string, len(allowlist))
	for _, key := range allowlist {
		if val, ok := os.LookupEnv(key); ok {
			result[key] = val
		}
	}
	return result
}
