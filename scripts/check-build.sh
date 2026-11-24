#!/bin/sh
set -eu

binary=$1
version=$2
commit=$3
commit_date=$4

metadata=$("$binary" version)
printf '%s\n' "$metadata" | grep -Fx "  version: $version"
printf '%s\n' "$metadata" | grep -Fx "  commit: $commit"

log=$(mktemp)
"$binary" server >"$log" 2>&1 &
pid=$!
cleanup() {
    status=$?
    kill "$pid" 2>/dev/null || true
    wait "$pid" 2>/dev/null || true
    if [ "$status" -ne 0 ]; then
        cat "$log" >&2
    fi
    rm -f "$log"
}
trap cleanup EXIT

body=$(curl -fsS --retry 5 --retry-connrefused --max-time 2 http://127.0.0.1:8080/version)
kill -0 "$pid"
printf '%s\n' "$body" | jq -e --arg version "$version" --arg commit "$commit" --arg date "$commit_date" \
    '.VERSION == $version and .COMMIT == $commit and .DATE == $date'

if output=$("$binary" server 2>&1); then
    echo 'Second server unexpectedly started successfully' >&2
    exit 1
fi
printf '%s\n' "$output" | grep -F 'serve HTTP on :8080:'
