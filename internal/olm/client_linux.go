//go:build linux

package olm

// CompanionAgentName is the OLM agent identifier used by the Praxis Linux Client.
// Must match the Agent value in the Linux app's OLM tunnel config.
const CompanionAgentName = "Praxis Linux"

// desktopAppSocketPath is the socket the Praxis Linux client's OLM listens
// on. It is distinct from defaultSocketPath, which the CLI's own client uses.
const desktopAppSocketPath = "/run/praxis-client/olm.sock"
