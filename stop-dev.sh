#!/usr/bin/env bash
set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PID_FILE="$SCRIPT_DIR/.api.pid"

if [ ! -f "$PID_FILE" ]; then
    echo "No PID file found. API server may not be running."
    exit 0
fi

PID=$(cat "$PID_FILE")
if kill -0 "$PID" 2>/dev/null; then
    echo "Stopping API server (PID: $PID)..."
    kill "$PID"
    rm -f "$PID_FILE"
    echo "API server stopped."
else
    echo "Process $PID not running. Cleaning up PID file."
    rm -f "$PID_FILE"
fi
