APP_NAME = rblog
DOCKER_COMPOSE = docker-compose

.PHONY: all
all: clean run

.PHONY: clean run

clean:
	@$(DOCKER_COMPOSE) down

run:
	@$(DOCKER_COMPOSE) up --build -d
	@docker rmi -f $$(docker images --filter "dangling=true" -q)