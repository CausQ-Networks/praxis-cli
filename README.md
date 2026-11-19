# Praxis CLI

Official CLI and client for [Praxis](https://praxis.causq.com) by CausQ Networks.

## Module

```
github.com/CausQ-Networks/praxis-cli
```

Upstream data-plane libraries remain `github.com/fosrl/newt` and `github.com/fosrl/olm` (they still read `PANGOLIN_ENDPOINT`; the CLI sets that alongside `PRAXIS_ENDPOINT`).

## Defaults

| Setting | Value |
| --- | --- |
| Cloud control plane | `https://app.praxis.causq.com` |
| Static CDN / installer | `https://praxis.causq.com/get-cli.sh` |
| Config dir | `~/.config/praxis` |
| Env prefix | `PRAXIS_CLI_*` |

Self-hosted: `praxis auth login --hostname https://your-instance` (or pick self-hosted in the prompt).

## Install

```bash
curl -fsSL https://praxis.causq.com/get-cli.sh | bash
```

## Build

```bash
make build
make go-build-release VERSION=0.21.2
```
