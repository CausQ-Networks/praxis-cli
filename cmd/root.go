package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/CausQ-Networks/praxis-cli/cmd/apply"
	"github.com/CausQ-Networks/praxis-cli/cmd/auth"
	"github.com/CausQ-Networks/praxis-cli/cmd/auth/login"
	"github.com/CausQ-Networks/praxis-cli/cmd/auth/logout"
	"github.com/CausQ-Networks/praxis-cli/cmd/authdaemon"
	companioncmd "github.com/CausQ-Networks/praxis-cli/cmd/companion"
	configcmd "github.com/CausQ-Networks/praxis-cli/cmd/config"
	"github.com/CausQ-Networks/praxis-cli/cmd/configure"
	"github.com/CausQ-Networks/praxis-cli/cmd/down"
	"github.com/CausQ-Networks/praxis-cli/cmd/list"
	"github.com/CausQ-Networks/praxis-cli/cmd/logs"
	"github.com/CausQ-Networks/praxis-cli/cmd/resetdns"
	"github.com/CausQ-Networks/praxis-cli/cmd/scp"
	selectcmd "github.com/CausQ-Networks/praxis-cli/cmd/select"
	"github.com/CausQ-Networks/praxis-cli/cmd/service"
	"github.com/CausQ-Networks/praxis-cli/cmd/ssh"
	"github.com/CausQ-Networks/praxis-cli/cmd/status"
	"github.com/CausQ-Networks/praxis-cli/cmd/up"
	"github.com/CausQ-Networks/praxis-cli/cmd/update"
	"github.com/CausQ-Networks/praxis-cli/cmd/version"
	"github.com/CausQ-Networks/praxis-cli/cmd/watchdog"
	"github.com/CausQ-Networks/praxis-cli/internal/api"
	"github.com/CausQ-Networks/praxis-cli/internal/companion"
	"github.com/CausQ-Networks/praxis-cli/internal/config"
	"github.com/CausQ-Networks/praxis-cli/internal/logger"
	"github.com/CausQ-Networks/praxis-cli/internal/notice"
	versionpkg "github.com/CausQ-Networks/praxis-cli/internal/version"
	"github.com/spf13/cobra"
)

// Initialize a root Cobra command.
//
// Set initResources to false when generating documentation to avoid
// parsing configuration files and instantiating the API client, among
// other such external resources. This is to avoid depending on external
// state when doing doc generation.
func RootCommand(initResources bool) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "praxis",
		Short:        "Praxis CLI",
		SilenceUsage: true,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		PersistentPreRunE: mainCommandPreRun,
	}

	cmd.AddCommand(auth.AuthCommand())
	if authDaemonCmd := authdaemon.AuthDaemonCmd(); authDaemonCmd != nil {
		cmd.AddCommand(authDaemonCmd)
	}
	cmd.AddCommand(apply.ApplyCommand())
	if companionCmd := companioncmd.CompanionCmd(); companionCmd != nil {
		cmd.AddCommand(companionCmd)
	}
	cmd.AddCommand(selectcmd.SelectCmd())
	cmd.AddCommand(list.ListCmd())
	cmd.AddCommand(configcmd.ConfigCmd())
	cmd.AddCommand(configure.ConfigureCmd())

	// Platform-specific commands - nil on unsupported platforms
	if upCmd := up.UpCmd(); upCmd != nil {
		cmd.AddCommand(upCmd)
	}
	if downCmd := down.DownCmd(); downCmd != nil {
		cmd.AddCommand(downCmd)
	}
	if logsCmd := logs.LogsCmd(); logsCmd != nil {
		cmd.AddCommand(logsCmd)
	}
	if statusCmd := status.StatusCmd(); statusCmd != nil {
		cmd.AddCommand(statusCmd)
	}
	if resetDNSCmd := resetdns.ResetDNSCmd(); resetDNSCmd != nil {
		cmd.AddCommand(resetDNSCmd)
	}
	if watchdogCmd := watchdog.WatchdogCmd(); watchdogCmd != nil {
		cmd.AddCommand(watchdogCmd)
	}

	cmd.AddCommand(service.ServiceCmd())
	cmd.AddCommand(ssh.SSHCmd())
	cmd.AddCommand(scp.SCPCmd())
	cmd.AddCommand(update.UpdateCmd())
	cmd.AddCommand(version.VersionCmd())
	cmd.AddCommand(login.LoginCmd())
	cmd.AddCommand(logout.LogoutCmd())

	if !initResources {
		return cmd, nil
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	logger.InitLogger(cfg.LogLevel)

	ctx := context.Background()
	ctx = config.WithConfig(ctx, cfg)
	cmd.SetContext(ctx)

	return cmd, nil
}

func commandNeedsAuthInit(cmd *cobra.Command) bool {
	// `up site` runs Newt standalone, exactly like the newt binary itself -
	// it takes its own id/secret/endpoint and never needs a Praxis
	// account or the companion daemon.
	if cmd.Name() == "site" && commandHasAncestor(cmd, "up") {
		return false
	}
	// `service` installs/manages systemd units and always takes its own
	// id/secret/endpoint (site) or id/secret (client) explicitly, so it
	// never needs a Praxis account or the companion daemon either.
	if cmd.Name() == "service" || commandHasAncestor(cmd, "service") {
		return false
	}
	for c := cmd; c != nil; c = c.Parent() {
		if c.Name() == "companion" || c.Name() == "config" {
			return false
		}
	}
	return true
}

func commandNeedsCompanionReady(cmd *cobra.Command) bool {
	if !commandNeedsAuthInit(cmd) {
		return false
	}
	return !commandExemptFromCompanionReady(cmd)
}

func commandExemptFromCompanionReady(cmd *cobra.Command) bool {
	switch cmd.Name() {
	case "version", "update", "login", "logout", "account", "org":
		return true
	}

	if cmd.Name() == "status" && commandHasAncestor(cmd, "auth") {
		return true
	}

	return false
}

func commandHasAncestor(cmd *cobra.Command, name string) bool {
	for p := cmd.Parent(); p != nil; p = p.Parent() {
		if p.Name() == name {
			return true
		}
	}
	return false
}

func initAuthContext(ctx context.Context, cfg *config.Config) (context.Context, error) {
	if _, ok := api.ClientFromContext(ctx); ok {
		return ctx, nil
	}

	companionState, accountStore, err := companion.Resolve(cfg)
	if err != nil {
		return ctx, err
	}

	var apiBaseURL string
	var sessionToken string

	if activeAccount, _ := accountStore.ActiveAccount(); activeAccount != nil {
		apiBaseURL = activeAccount.Host
		sessionToken = activeAccount.SessionToken
	}

	client, err := api.InitClient(apiBaseURL, sessionToken, cfg.SessionCookieName)
	if err != nil {
		return ctx, err
	}

	ctx = api.WithAPIClient(ctx, client)
	ctx = config.WithAccountStore(ctx, accountStore)
	ctx = companion.WithState(ctx, companionState)
	return ctx, nil
}

func mainCommandPreRun(cmd *cobra.Command, args []string) error {
	cfg := config.ConfigFromContext(cmd.Context())
	if cfg == nil {
		return fmt.Errorf("configuration not loaded")
	}

	if err := notice.ShowPending(cfg); err != nil {
		logger.Debug("Failed to show pending notices: %v", err)
	}

	if commandHasAncestor(cmd, "config") {
		// Config commands must work without an account (or with the desktop
		// app not running), but use the active org when there is one.
		if ctx, err := initAuthContext(cmd.Context(), cfg); err == nil {
			cmd.SetContext(ctx)
		} else {
			logger.Debug("Continuing without account for %q: %v", cmd.Name(), err)
		}
	} else if commandNeedsAuthInit(cmd) {
		ctx, err := initAuthContext(cmd.Context(), cfg)
		if err != nil {
			return err
		}
		cmd.SetContext(ctx)
		cfg = config.ConfigFromContext(ctx)

		if commandNeedsCompanionReady(cmd) {
			if err := companion.GuardReady(ctx); err != nil {
				return err
			}
		}
	}

	// Skip update checks when running self-update or companion commands.
	cmdName := cmd.Name()
	if cmdName == "update" || !commandNeedsAuthInit(cmd) {
		logger.Debug("Skipping update check for %q command", cmdName)
		return nil
	}

	ensureRuntimeDirs(cfg)

	// Check for updates asynchronously
	if !cfg.DisableUpdateCheck {
		logger.Debug("Starting update check for %q command", cmdName)
		versionpkg.CheckForUpdateAsync(func(release *versionpkg.GitHubRelease) {
			logger.Warning("A new version is available: %s (current: %s)", release.TagName, versionpkg.Version)
			logger.Info("Run 'praxis update' to update to the latest version")
			fmt.Println()
		})
	} else {
		logger.Debug("Update check disabled by configuration")
	}

	return nil
}

// Make sure all required directories exist once
// before executing any subcommands.
func ensureRuntimeDirs(cfg *config.Config) {
	configDir, err := config.GetPraxisConfigDir()
	if err != nil {
		logger.Warning("failed to create praxis configuration directory: %v", err)
	} else {
		err = os.MkdirAll(configDir, 0o755)
		if err != nil {
			logger.Warning("failed to create %s: %v", configDir, err)
		}
	}

	if cfg.LogFile != "" {
		logPathDirname := filepath.Dir(cfg.LogFile)

		err = os.MkdirAll(logPathDirname, 0o755)
		if err != nil {
			logger.Warning("failed to create %s: %v", logPathDirname, err)
		}
	}
}

// Execute is called by main.go
func Execute() {
	cmd, err := RootCommand(true)
	if err != nil {
		logger.Error("%v", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
