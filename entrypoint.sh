#!/bin/sh

set -e

# Call this as the fist thing so that if the user has added custom certificates
# to /usr/local/share/ca-certificates/ they will be used by praxis-cli.
update-ca-certificates

# Check if Praxis environment variables are set
if [ -n "$PRAXIS_ENDPOINT" ] && [ -n "$CLIENT_ID" ] && [ -n "$CLIENT_SECRET" ]; then
    # Run praxis-cli up --attach with the provided credentials
    exec praxis-cli up client --attach --id "$CLIENT_ID" --secret "$CLIENT_SECRET" --endpoint "$PRAXIS_ENDPOINT" "$@"
fi

# Check if site (Newt) environment variables are set. SITE_ID/SITE_SECRET
# (along with PRAXIS_ENDPOINT and any other newt flags/env vars) are read
# directly by `up site`, so nothing needs to be translated into flags here.
if [ -n "$SITE_ID" ] && [ -n "$SITE_SECRET" ]; then
    exec praxis-cli up site "$@"
fi

# If no arguments provided, run praxis-cli with default behavior
if [ $# -eq 0 ]; then
    exec praxis-cli
fi

# If first arg is a flag (starts with -), prepend praxis-cli
if [ "${1#-}" != "$1" ]; then
    exec praxis-cli "$@"
fi

# If first arg is a shell or absolute path, execute as-is
case "$1" in
    sh|bash|/*)
        exec "$@"
        ;;
    *)
        # Otherwise, treat as a praxis-cli subcommand
        exec praxis-cli "$@"
        ;;
esac