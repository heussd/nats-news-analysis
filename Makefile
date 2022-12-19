SHELL             = /bin/bash
CPU_CORES         = $(shell nproc)
SCALE_PERFORMANCE = $$(($(CPU_CORES)*70/100))
SCALE_POWER_SAFE  = $$(($(CPU_CORES)*30/100))

all:
	@echo "Scaling to $(SCALE_PERFORMANCE) (performance)"
	docker-compose up -d \
		--scale keyword-matcher-go=$(SCALE_PERFORMANCE) \
		--scale keyword-matcher-python=0 \
		--scale rss-article-url-feeder=3

build:
	docker-compose build

logs:
	docker-compose logs -f keyword-matcher-go

logs-feeder:
	docker-compose logs -f rss-article-url-feeder


stop:
	docker-compose down

watch:
	watch nats stream ls

power-safe-mode:
	@echo "Scaling to $(SCALE_POWER_SAFE) (power safe)"
	docker-compose up -d \
		--scale keyword-matcher-go=$(SCALE_POWER_SAFE) \
		--scale keyword-matcher-python=0 \


