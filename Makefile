SHELL             = /bin/bash
CPU_CORES         = $(shell nproc)
SCALE_PERFORMANCE = $$(($(CPU_CORES)*70/100))
SCALE_POWER_SAFE  = $$(($(CPU_CORES)*30/100))


kitty:
	kitty @launch --location split --cwd=current make watch
	kitty @launch --location split --cwd=current ctop 
	$(MAKE) run


run: start
	bash -c "trap 'trap - SIGINT SIGTERM ERR; $(MAKE) stop; exit 1' SIGINT SIGTERM ERR; $(MAKE) logs"


build:
	docker-compose \
		-f docker-compose.build.yml \
		-f docker-compose.override.yml \
		build


logs:
	docker-compose logs -f keyword-matcher-go


stop:
	docker-compose down


watch:
	watch nats stream ls


start:
	docker-compose up -d 


performance-mode:
	@echo "Scaling to $(SCALE_PERFORMANCE) (performance)"
	docker-compose up -d \
		--scale keyword-matcher-go=$(SCALE_PERFORMANCE) \
		--scale keyword-matcher-python=0 \
		--scale rss-article-url-feeder=3


# https://yuriktech.com/2020/03/21/Collecting-Docker-Logs-With-Loki/
install-loki-driver:
	docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions