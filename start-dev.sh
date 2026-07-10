#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
API_DIR="$SCRIPT_DIR/api"
PID_FILE="$SCRIPT_DIR/.api.pid"

if [ -f "$PID_FILE" ]; then
    echo "API server may already be running (PID file exists). Run stop-dev.sh first."
    exit 1
fi

echo "Starting POS API server..."
cd "$API_DIR"

# Build
go build -o pos-api ./cmd/server

# Run in background
./pos-api &
API_PID=$!
echo "$API_PID" > "$PID_FILE"

echo "API server started (PID: $API_PID)"
echo "Listening on http://localhost:${SERVER_PORT:-8080}"
