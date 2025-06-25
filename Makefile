reup:
	docker-compose down
	docker-compose up -d --build

up:
	docker-compose up -d --build

down:
	docker-compose down

dump:
	docker-compose run --rm backup