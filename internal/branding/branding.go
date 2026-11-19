// Package branding holds Praxis product defaults shared across the CLI.
package branding

const (
	// ProductName is the user-facing product name.
	ProductName = "Praxis"

	// BinaryName is the installed CLI command name.
	BinaryName = "praxis"

	// ModulePath is the Go module path for this repository.
	ModulePath = "github.com/CausQ-Networks/praxis-cli"

	// GitHubRepo is owner/name for releases and update checks.
	GitHubRepo = "CausQ-Networks/praxis-cli"

	// CloudHost is the default Praxis Cloud control-plane hostname (no scheme).
	CloudHost = "app-praxis.causq.com"

	// CloudURL is the default Praxis Cloud base URL (dashboard / API root).
	CloudURL = "https://" + CloudHost

	// StaticCDNURL hosts installers and static assets.
	StaticCDNURL = "https://praxis.causq.com"

	// InstallerURL is the curl|bash installer script.
	InstallerURL = StaticCDNURL + "/get-cli.sh"

	// ReleasesAPIURL is the GitHub latest-release endpoint.
	ReleasesAPIURL = "https://api.github.com/repos/" + GitHubRepo + "/releases/latest"

	// ConfigDirName is the directory under ~/.config for CLI state.
	ConfigDirName = "praxis"

	// EnvPrefix is the viper/env prefix for CLI settings (PRAXIS_CLI_*).
	EnvPrefix = "PRAXIS_CLI"

	// EndpointEnv is the preferred env var for the control-plane URL.
	EndpointEnv = "PRAXIS_ENDPOINT"

	// LegacyEndpointEnv is still read by upstream newt/olm; set alongside EndpointEnv.
	LegacyEndpointEnv = "PANGOLIN_ENDPOINT"
)
