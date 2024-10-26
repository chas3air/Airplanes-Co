build:
	@docker compose build

up:
	@docker compose up

down:
	@docker compose down

refresh:
	@docker compose down
	@docker compose build
	@docker compose up