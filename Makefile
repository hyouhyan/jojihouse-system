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

dev/up:
	docker compose -f dev-compose.yml up -d --build

dev/reup:
	make dev/down
	make dev/up

dev/down:
	docker compose -f dev-compose.yml down
