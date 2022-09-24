SHELL             = /bin/bash
CPU_CORES         = $(shell nproc)
SCALE_PERFORMANCE = $$(($(CPU_CORES)*70/100))
SCALE_POWER_SAFE  = $$(($(CPU_CORES)*30/100))

all:
	@echo "Scaling to $(SCALE_PERFORMANCE) (performance)"
	docker-compose up -d \
		--scale nats-news-keyword-matcher=$(SCALE_PERFORMANCE) \
		--scale nats-rss-article-url-feeder=3

stop:
	docker-compose down
watch:
	watch nats stream info article-urls
retrieve-only:
	docker-compose up -d \
		--scale nats-news-keyword-matcher=0 \
		--scale fullfeedrss=0 \
		--scale nats-pocket-integration=0 \
		--scale nats-rss-article-url-feeder=3

power-safe-mode:
	@echo "Scaling to $(SCALE_POWER_SAFE) (power safe)"
	docker-compose up -d \
		--scale nats-news-keyword-matcher=$(SCALE_POWER_SAFE)

