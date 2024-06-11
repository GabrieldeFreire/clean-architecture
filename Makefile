ENV_FILE = cmd/ordersystem/.env

include $(ENV_FILE)
$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' $(ENV_FILE)))

PROJECT_NAME := clean-app
VERSION := latest

COMPOSE_RUN := docker compose up -V
DOCKER_BUILD := docker build --no-cache -t $(PROJECT_NAME):$(VERSION) .
COMPOSE_DOWN := docker compose down --volumes


build:
	$(DOCKER_BUILD)

run: build
	bash -c "trap '$(COMPOSE_DOWN)' EXIT; $(COMPOSE_RUN)"