.SILENT:

build:
	go build -o hzpaste-api.bin cmd/hzpaste-api/main.go

test:
	HZPASTE_PORT=8889 HZPASTE_HOST=0.0.0.0 go test -v ./cmd/hzpaste-api/

run:
	HZPASTE_PORT=8888 HZPASTE_HOST=0.0.0.0 go run cmd/hzpaste-api/main.go

release: build
	HZPASTE_PORT=8888 HZPASTE_HOST=0.0.0.0 GIN_MODE=release ./hzpaste-api

swag-init:
	HZPASTE_PORT=8888 HZPASTE_HOST=0.0.0.0 cd cmd/hzpaste-api/ && swag init -g main.go -o internal/docs/ --parseDependency --parseInternal
