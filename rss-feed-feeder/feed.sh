#!/bin/bash
set -o errexit
set -o pipefail
set -o nounset

sleep 5

nats stream add feed-urls \
    --server="$NATS_SERVER" \
    --subjects feed-url \
    --storage file \
    --replicas 1 \
    --retention work \
    --discard old \
    --max-msgs=-1 \
    --max-msg-size=-1 \
    --max-msgs-per-subject=-1 \
    --max-bytes=-1 \
    --max-age=-1  \
    --allow-rollup \
    --allow-direct \
    --no-deny-delete \
    --no-deny-purge \
    --dupe-window=10m0s 2> /dev/null || true


nats stream info feed-urls \
    --server="$NATS_SERVER"

while true; do
    while read feed; do
        case "$feed" in
            \#*)
                ;;
            *)
                nats pub feed-url "$feed"  \
                    --server="$NATS_SERVER" \
                    --header="Nats-Msg-Id:$feed" # De-duplication works with Msg Id
        esac
    done </urls.txt

    echo "Waiting 2h..."
    sleep 7200
done