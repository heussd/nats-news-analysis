#!/bin/bash
set -o errexit
set -o pipefail
set -o nounset

sleep 5


# Stream config needs to be exactly in sync with go config.
# Create via go and export via $nats stream info --json 
./nats stream create \
    --server="$NATS_SERVER" \
    --config feed-urls.json 


./nats stream info feed-urls \
    --server="$NATS_SERVER"

while true; do
    while read feed; do
        case "$feed" in
            \#*)
                ;;
            *)
                ./nats pub feed-urls "{\"url\": \"$feed\"}"  \
                    --server="$NATS_SERVER" \
                    --header="Nats-Msg-Id:$feed" # De-duplication works with Msg Id
        esac
    done </urls.txt

    echo "Waiting 2h..."
    sleep 7200
done