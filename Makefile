all:
	docker-compose up -d
stop:
	docker-compose down
watch:
	watch nats stream info article-urls
retrieve-only:
	docker-compose up -d \
		--scale nats-news-keyword-matcher=0 \
		--scale fullfeedrss=0 \
		--scale nats-pocket-integration=0
