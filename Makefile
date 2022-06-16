VERSION=$(shell git rev-parse --short HEAD)

build-dev: fmt-code vet-code
	docker-compose build --build-arg APP_VERSION=$(VERSION)

up-dev:
	docker-compose up server

# Gofmt is a tool that automatically formats Go source code.
fmt-code:
	go fmt ./...

# Vet examines Go source code and reports suspicious constructs,
# such as Printf calls whose arguments do not align with the format string.
vet-code:
	go vet ./...