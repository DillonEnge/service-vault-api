reset:
	docker compose down && docker compose up -d --build

down:
	docker-compose down

dev:
	docker-compose up -d

build:
	docker-compose up -d --build
