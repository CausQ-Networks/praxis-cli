package scp

import sshcmd "github.com/CausQ-Networks/praxis-cli/cmd/ssh"

type RunOpts struct {
	User          string
	Hostname      string
	Port          int
	PrivateKeyPEM string
	Certificate   string
	ResourceID    string
	Passthrough   sshcmd.SSHPassthrough
}
