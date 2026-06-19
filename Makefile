.PHONY: api web build-api build-web test fmt

api:
	cd api && go run ./cmd/server

web:
	cd web && npm run dev

build-api:
	cd api && go build -o bin/server ./cmd/server

build-web:
	cd web && npm run build

test:
	cd api && go test ./...

fmt:
	cd api && gofmt -w .
