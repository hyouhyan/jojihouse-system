.PHONY: reup up down clean dump deploy dev dev/up dev/down dev/reup

reup:
	make down
	make up

up:
	docker compose up -d --build
	make clean

down:
	docker compose down

clean:
	docker image prune -f

deploy:
	git fetch
	make down
	git pull
	make up

dump:
	docker compose run --rm backup

dev:
	make dev/reup
	@echo "Dev mode is \033[32mactive\033[m."
	@echo "Try this -> \033[4m\033[34mhttp://localhost:1024/\033[m\033[m"

dev/up:
	docker compose -f dev-compose.yml up -d --build

dev/reup:
	make dev/down
	make dev/up

dev/down:
	docker compose -f dev-compose.yml down
