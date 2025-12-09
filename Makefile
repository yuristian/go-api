DB_URL ?= $(DATABASE_URL)

migrate-create:
	go run ./cmd/migrate -action create -name "$(name)"

migrate-up:
	go run ./cmd/migrate -action up

migrate-down:
	go run ./cmd/migrate -action down -steps 1

migrate-version:
	go run ./cmd/migrate -action version

migrate-force:
	go run ./cmd/migrate -action force -version $(version)
