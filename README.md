# Praxis CLI

[![Discord](https://img.shields.io/discord/1572896060949659679?logo=discord&style=flat-square&label=Discord)](https://praxis.causq.com/discord/)

Official CLI and client for [Praxis](https://praxis.causq.com) by CausQ Networks.

## Module

```
github.com/CausQ-Networks/praxis-cli
```

Data-plane libraries are `github.com/CausQ-Networks/praxis-newt` and `github.com/CausQ-Networks/praxis-olm` (CLI sets `PRAXIS_ENDPOINT` and mirrors `PANGOLIN_ENDPOINT` for compatibility).

## Defaults

| Setting | Value |
| --- | --- |
| Cloud control plane | `https://praxis.causq.com` (app.praxis.causq.com once DNS is attached) |
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
