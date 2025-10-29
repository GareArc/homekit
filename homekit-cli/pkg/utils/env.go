package utils

import (
	"os"
	"strings"
)

// MergeEnv merges two maps giving precedence to overrides.
func MergeEnv(base map[string]string, overrides map[string]string) map[string]string {
	out := make(map[string]string, len(base)+len(overrides))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overrides {
		out[k] = v
	}
	return out
}

// FilterInheritedEnv filters environment variables by allowed prefixes.
func FilterInheritedEnv(prefixes []string) []string {
	var result []string
nextVar:
	for _, kv := range os.Environ() {
		for _, prefix := range prefixes {
			if strings.HasPrefix(kv, prefix) {
				result = append(result, kv)
				continue nextVar
			}
		}
	}
	return result
}

// SplitPair splits a string into key/value using the first separator occurrence.
func SplitPair(value string, sep string) (parts [2]string) {
	idx := strings.Index(value, sep)
	if idx < 0 {
		parts[0] = strings.TrimSpace(value)
		return
	}
	parts[0] = strings.TrimSpace(value[:idx])
	parts[1] = strings.TrimSpace(value[idx+len(sep):])
	return
}
