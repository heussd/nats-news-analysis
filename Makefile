SHELL             = /bin/bash
CPU_CORES         = $(shell nproc)
SCALE_PERFORMANCE = $$(($(CPU_CORES)*70/100))
SCALE_POWER_SAFE  = $$(($(CPU_CORES)*30/100))


start:
	docker compose up --build -d
	open http://localhost:3000/d/QyuE2Of4z/news-analysis?orgId=1&refresh=5s


logs:
	docker compose logs -f keyword-matcher-go


stop:
	docker compose down


watch:
	watch nats stream ls


# https://yuriktech.com/2020/03/21/Collecting-Docker-Logs-With-Loki/
install-loki-driver:
	docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions
