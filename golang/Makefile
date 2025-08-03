publish:
	nats pub article-url "https://www.tagesschau.de/"


watch-articles:
	 watch nats stream info article-urls


watch-matches:
	 watch nats stream info match_urls


stream:
	nats stream add article-urls --subjects article-url --storage memory --replicas 1 --retention work --discard old --max-msgs=-1 --max-msgs-per-subject=-1 --max-bytes=-1 --max-age=-1  --allow-rollup --allow-direct


bake:
	COMMIT_REF=$(shell git rev-parse --short HEAD ) \
		docker bake --push