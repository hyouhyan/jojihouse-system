DOCKER=docker
COMPOSE=$(DOCKER) compose

IMAGE=$(DOCKER) image

.PHONY: reup
reup:
	$(MAKE) down
	$(MAKE) up

.PHONY:  up
up:
	$(COMPOSE) up -d --build
	$(MAKE) clean

.PHONY:  down
down:
	$(COMPOSE) down

.PHONY:  clean
clean:
	$(IMAGE) prune -f

.PHONY:  deploy
deploy:
	git fetch
	@if [ $$(git rev-parse HEAD) = $$(git rev-parse @{u}) ]; then \
		echo "No updates. Skip deploy."; \
	else \
		echo "Updates found. Deploying..."; \
		$(MAKE) down; \
		git pull; \
		$(MAKE) up; \
	fi


.PHONY:  dump
dump:
	$(COMPOSE) run --rm backup

.PHONY:  dev
dev:
	$(MAKE) dev/reup
	@echo "Dev mode is \033[32mactive\033[m."
	@echo "Try this -> \033[4m\033[34mhttp://localhost:1024/\033[m\033[m"

.PHONY:  dev/up
dev/up:
	$(COMPOSE) -f dev-compose.yml up -d --build

.PHONY:  dev/down
dev/reup:
	$(MAKE) dev/down
	$(MAKE) dev/up

.PHONY:  dev/reup
dev/down:
	$(COMPOSE) -f dev-compose.yml down
