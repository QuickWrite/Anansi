#!/bin/sh
set -e

: "${ANANSI_LIMIT:=1000}"
: "${ANANSI_PORT:=8080}"
: "${ANANSI_APPLICATION_NAME:=Anansi}"
: "${ANANSI_TITLE_SEPARATOR:=-}"
: "${ANANSI_TITLE_POSITION:=left}"
: "${ANANSI_LINK_PROBABILITY:=5}"
: "${ANANSI_MIN_LINK_LENGTH:=4}"

if [ -n "${ANANSI_INPUT_DIR:-}" ] && [ -n "${ANANSI_PATH:-}" ]; then
    echo "ERROR: ANANSI_PATH and ANANSI_INPUT_DIR cannot be set simultaneously." >&2
    exit 1
fi

if [ -n "${ANANSI_INPUT_DIR:-}" ]; then
    set -- "$ANANSI_INPUT_DIR"/*
else
    : "${ANANSI_PATH:=./default.txt}"
    set -- "$ANANSI_PATH"
fi

exec ./Anansi "$@" \
             --limit="${ANANSI_LIMIT}" \
             --port="${ANANSI_PORT}" \
             --title="${ANANSI_APPLICATION_NAME},${ANANSI_TITLE_SEPARATOR},${ANANSI_TITLE_POSITION}" \
             --r="${ANANSI_LINK_PROBABILITY}" \
             --w="${ANANSI_MIN_LINK_LENGTH}"
