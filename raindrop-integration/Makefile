publish:
	nats pub match-url "{\"RegexId\": \"OMG\", \"Url\": \"https://www.tagesschau.de/ausland/europa/afrika-russland-ukraine-frieden-100.html\"}"

stream:
	nats stream add match-urls --subjects match-url --storage memory --replicas 1 --retention work --discard old --max-msgs=-1 --max-msgs-per-subject=-1 --max-bytes=-1 --max-age=-1  --allow-rollup --allow-direct
