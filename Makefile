.SILENT:

build:
	go build cmd/hzpaste-api/main.go

run:
	go run cmd/hzpaste-api/main.go

release:
	GIN_MODE=release go run cmd/hzpaste-api/main.go

