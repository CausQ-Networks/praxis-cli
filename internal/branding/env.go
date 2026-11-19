package branding

import "os"

// EndpointEnvVars returns env vars to publish a control-plane URL so both
// Praxis CLI and upstream newt/olm (which still read PANGOLIN_ENDPOINT) see it.
func EndpointEnvVars(endpoint string) map[string]string {
	if endpoint == "" {
		return nil
	}
	return map[string]string{
		EndpointEnv:       endpoint,
		LegacyEndpointEnv: endpoint,
	}
}

// FirstEnv returns the first non-empty environment value among names.
func FirstEnv(names ...string) string {
	for _, name := range names {
		if v := os.Getenv(name); v != "" {
			return v
		}
	}
	return ""
}
