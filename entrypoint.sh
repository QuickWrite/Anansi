#!/bin/sh
set -e

: "${ANANSI_PATH:=./default.txt}"
: "${ANANSI_LIMIT:=1000}"
: "${ANANSI_PORT:=8080}"
: "${ANANSI_APPLICATION_NAME:=Anansi}"
: "${ANANSI_TITLE_SEPARATOR:=-}"
: "${ANANSI_TITLE_POSITION:=left}"
: "${ANANSI_LINK_PROBABILITY:=5}"
: "${ANANSI_MIN_LINK_LENGTH:=4}"

exec ./Anansi "${ANANSI_PATH}" \
             --limit="${ANANSI_LIMIT}" \
             --port="${ANANSI_PORT}" \
             --title="${ANANSI_APPLICATION_NAME},${ANANSI_TITLE_SEPARATOR},${ANANSI_TITLE_POSITION}" \
             --r="${ANANSI_LINK_PROBABILITY}" \
             --w="${ANANSI_MIN_LINK_LENGTH}"
