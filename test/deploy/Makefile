.PHONY: default
default: provision start

.PHONY: provision
provision:
	docker compose -f postgresql/docker-compose.yaml up -d

.PHONY: start
start:
	cd postgresql && ./configure-postgresql.sh

.PHONY: deprovision
deprovision:
	docker compose -f postgresql/docker-compose.yaml down