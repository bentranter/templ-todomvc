all: run

run:
	go generate ./... && go run cmd/main.go

.PHONY: run
