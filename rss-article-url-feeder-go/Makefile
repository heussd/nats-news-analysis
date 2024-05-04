stream:
	nats stream add feed-urls \
		--subjects feed-url \
		--storage memory \
		--replicas 1 \
		--retention work \
		--discard old \
		--max-msgs=-1 \
		--max-msgs-per-subject=-1 \
		--max-bytes=-1 \
		--max-age=-1 \
		--allow-rollup \
		--allow-direct \
		--max-msg-size=-1 \
		--dupe-window=2m0s \
		--no-deny-delete \
		--no-deny-purge

