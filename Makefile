.PHONY: reup up down clean dump

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
	git pull
	make reup

dump:
	docker compose run --rm backup
