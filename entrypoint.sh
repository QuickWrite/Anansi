#!/bin/sh
set -e

: "${ANANSI_PATH:=./default.txt}"
: "${ANANSI_LIMIT:=1000}"
: "${ANANSI_PORT:=8080}"
: "${ANANSI_LINK_PROBABILITY:=5}"
: "${ANANSI_MIN_LINK_LENGTH:=4}"

exec ./Anansi "${ANANSI_PATH}" \
             --limit="${ANANSI_LIMIT}" \
             --port="${ANANSI_PORT}" \
             --r="${ANANSI_LINK_PROBABILITY}" \
             --w="${ANANSI_MIN_LINK_LENGTH}"
