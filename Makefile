local:
	go run ./engine/rest/rest.go

wire:
	wire ./src

migration-create:
	@read -p "Enter the migration name: " name; \
	go run ./engine/goose/goose.go -schema=migrations create $$name

migration-up:
	go run ./engine/goose/goose.go -schema=migrations up

migration-down:
	go run ./engine/goose/goose.go -schema=migrations down