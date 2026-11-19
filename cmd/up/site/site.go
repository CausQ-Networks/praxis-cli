// Package site implements `praxis up site`, which runs a Newt site
// tunnel embedded directly in the Praxis CLI. It accepts exactly the same
// command-line flags and environment variables as the standalone newt
// binary (see github.com/CausQ-Networks/praxis-newt), reusing newt's own config-loading and
// runtime packages rather than re-implementing them.
package site

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/CausQ-Networks/praxis-cli/internal/config"
	versionpkg "github.com/CausQ-Networks/praxis-cli/internal/version"
	"github.com/CausQ-Networks/praxis-newt/clients/permissions"
	newtLogger "github.com/CausQ-Networks/praxis-newt/logger"
	newtpkg "github.com/CausQ-Networks/praxis-newt/newt"
	"github.com/CausQ-Networks/praxis-newt/newtconfig"
	"github.com/CausQ-Networks/praxis-newt/updates"
	"github.com/CausQ-Networks/praxis-newt/websocket"
	"github.com/spf13/cobra"
)

// SiteUpCmd returns the `up site` command. Flag parsing is delegated
// entirely to newtconfig (the same flag set the newt binary itself parses),
// so this command intentionally does not declare any cobra flags of its
// own - run `praxis up site --help` to see them.
func SiteUpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "site",
		Short: "Start a site connection (Newt)",
		Long: `Bring up a site tunnel using Newt, embedded directly in the Praxis CLI.

This accepts the same flags and environment variables as the standalone
newt binary. Run 'praxis up site --help' to see them.`,
		DisableFlagParsing: true,
		SilenceUsage:       true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd.Context(), args)
		},
	}

	return cmd
}

func run(ctx context.Context, args []string) error {
	newtLogger.Init(nil)

	defaultConfigFile, err := defaultSiteConfigFile()
	if err != nil {
		return fmt.Errorf("failed to resolve default config file location: %w", err)
	}

	cfg, err := newtconfig.Load(newtconfig.Options{
		Args:              args,
		Version:           versionpkg.NewtVersion(),
		Agent:             "cli",
		AgentVersion:      versionpkg.Version,
		Platform:          runtime.GOOS,
		DefaultConfigFile: defaultConfigFile,
	})
	if err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	if cfg.UseNativeMainInterface {
		if err := permissions.CheckNativeInterfacePermissions(); err != nil {
			return fmt.Errorf("insufficient permissions for native main tunnel interface: %w", err)
		}
	}

	// Restart in place (preserving the original "up site" arguments)
	// rather than exiting, matching the standalone newt binary's behavior
	// when a blueprint reload requires a fresh tunnel.
	cfg.OnRestart = reexec

	sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	n, err := newtpkg.Init(sigCtx, cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize newt: %w", err)
	}

	resolvedCfg := n.GetConfig()

	startSelfUpdateChecks(sigCtx, resolvedCfg)

	n.Start(sigCtx)

	return nil
}

// defaultSiteConfigFile returns the site tunnel's config file path when run
// under the Praxis CLI: alongside the CLI's own config.json, in the same
// ~/.config/praxis directory, but named site.json so the two don't
// collide. This only applies as a fallback default - --config-file and
// CONFIG_FILE still take precedence, and the standalone newt binary keeps
// using its own default (~/.config/newt-client/config.json) since it never
// sets newtconfig.Options.DefaultConfigFile.
func defaultSiteConfigFile() (string, error) {
	dir, err := config.GetPraxisConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "site.json"), nil
}

// startSelfUpdateChecks periodically checks the Praxis server for a newer
// CLI release and, if one is available, downloads it, replaces the running
// `pangolin` binary on disk, and re-execs in place - mirroring the
// self-update loop the standalone newt binary runs for itself, but with
// Agent: "cli" so the server hands back a praxis-cli release instead of a
// newt one.
func startSelfUpdateChecks(ctx context.Context, resolvedCfg newtpkg.Config) {
	// Reuse the same TLS parameters as the websocket client for the
	// self-update HTTP requests.
	var selfUpdateTLS *tls.Config
	if resolvedCfg.TLSClientCert != "" || resolvedCfg.TLSPrivateKey != "" {
		selfUpdateTLS, _ = websocket.BuildTLSConfig(
			resolvedCfg.TLSClientCert,
			resolvedCfg.TLSClientKey,
			resolvedCfg.TLSClientCAs,
			resolvedCfg.TLSPrivateKey,
		)
	}

	doUpdate := func() {
		newtLogger.Debug("checkAndSelfUpdate: running periodic update check")
		if err := updates.CheckAndSelfUpdate(updates.SelfUpdateConfig{
			Endpoint:       resolvedCfg.Endpoint,
			NewtID:         resolvedCfg.ID,
			Secret:         resolvedCfg.Secret,
			CurrentVersion: versionpkg.NewtVersion(),
			TLSConfig:      selfUpdateTLS,
			Agent:          "cli",
			Restart:        updateRestart(),
		}); err != nil {
			if errors.Is(err, updates.ErrAutoUpdateUnsupportedInOfficialContainer) {
				newtLogger.Debug("checkAndSelfUpdate: auto-update skipped: %v", err)
				return
			}
			newtLogger.Error("Auto-update check failed: %v", err)
		}
	}

	go func() {
		time.Sleep(2 * time.Minute)
		doUpdate()
		ticker := time.NewTicker(6 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				doUpdate()
			case <-ctx.Done():
				return
			}
		}
	}()
}
