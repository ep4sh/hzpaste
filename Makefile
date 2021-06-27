.SILENT:

build:
	go build cmd/hzpaste-api/main.go

run:
	go run cmd/hzpaste-api/main.go

release:
	GIN_MODE=release go run cmd/hzpaste-api/main.go

swag-init:
	cd cmd/hzpaste-api/ && swag init -g main.go -o internal/docs/ --parseDependency --parseInternal
