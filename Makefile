reup:
	docker compose down
	docker compose up -d --build
	docker image prune -f

up:
	docker compose up -d --build
	docker image prune -f

down:
	docker compose down

dump:
	docker compose run --rm backup