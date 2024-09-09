ifneq (,$(wildcard .env))
    include .env
    export
endif

build:
	go build main.go

run: build
	./main

test:
	cd service/ && go test -v

test_bench:
	cd service/ && go test -bench=.

migrate_up:
	goose -dir migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" up

migrate_down:
	goose -dir migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" down

migrate_status:
	goose -dir migrations postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable" status

clean:
	rm -rf main
