.SILENT:

build:
	go build -o hzpaste-api.bin cmd/hzpaste-api/main.go

test:
	go test -v ./cmd/hzpaste-api/

run:
	go run cmd/hzpaste-api/main.go

release: build
	GIN_MODE=release ./hzpaste-api

swag-init:
	cd cmd/hzpaste-api/ && swag init -g main.go -o internal/docs/ --parseDependency --parseInternal
