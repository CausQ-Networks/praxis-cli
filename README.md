# Praxis CLI

Official CLI and client for [Praxis](https://praxis.causq.com) by CausQ Networks.

Forked from `fosrl/cli` and rebuilt with Praxis branding. Go module path remains `github.com/fosrl/cli` for compatibility with `newt` / `olm`.

## Install

```bash
curl -fsSL https://praxis.causq.com/get-cli.sh | bash
```

Binaries: [CausQ-Networks/praxis-cli releases](https://github.com/CausQ-Networks/praxis-cli/releases)

## Build

```bash
make build
# or multi-arch release artifacts:
make go-build-release VERSION=0.21.1
```
