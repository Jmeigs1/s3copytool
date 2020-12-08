package main

import "strings"

// Helper functions
func joinPrefixes(prefixes []string) string {
	if len(prefixes) == 0 {
		return ""
	}
	return strings.Join(prefixes, "/") + "/"
}
