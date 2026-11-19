package utils

import (
	"strings"

	"github.com/CausQ-Networks/praxis-cli/internal/branding"
)

const defaultHostname = branding.CloudHost

// FormatHostnameBaseURL returns the hostname formatted as a base URL (with protocol, without /api/v1).
// This is useful for constructing URLs to the web interface.
func FormatHostnameBaseURL(hostname string) string {
	if hostname == "" {
		hostname = defaultHostname
	}

	// Ensure hostname has protocol
	if !strings.HasPrefix(hostname, "http://") && !strings.HasPrefix(hostname, "https://") {
		hostname = "https://" + hostname
	}

	// Remove /api/v1 suffix if present
	hostname = strings.TrimSuffix(hostname, "/api/v1")
	hostname = strings.TrimSuffix(hostname, "/")

	return hostname
}
