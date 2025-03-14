.PHONY: up down build migrate logs swag

 up:
	docker-compose up --build

 down:
	docker-compose down

 build:
	docker-compose build

 migrate:
	docker-compose run --rm migrations

 logs:
	docker-compose logs -f app

 swag:
	swag init -g cmd/app/main.go