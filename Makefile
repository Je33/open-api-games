.PHONY: build build-prof build-run dc gen test run dev lint

build:
	go build -o ./build/api ./cmd/api/main.go

build-prof: build
	go tool pprof —text ./build/api

build-run: build
	./build/api

dc:
	docker-compose up  --remove-orphans --build

gen:
	go generate ./...

test:
	go test -v -coverprofile cover.out ./... && go tool cover -html=cover.out

run:
	go run -race ./cmd/api/main.go

dev:
	air -c api.air.toml

lint:
	golangci-lint run