# Makefile
migrate-new:
	@read -p "Enter migration name: " name; \
	dbmate new $$name

migrate-up:
	dbmate up

migrate-down:
	dbmate down

gen:
	sqlc generate
	swag init -g cmd/api/main.go