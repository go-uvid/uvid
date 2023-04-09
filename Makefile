lint:
	golangci-lint run -c ./golangci.yml ./...

test:
	go test ./... -v --cover
	cd portal && pnpm --filter "uvid-js" test

test-report:
	go test ./... -v --cover -coverprofile=coverage.out
	go tool cover -html=coverage.out

build:
	GOFLAGS=-mod=mod go build -o bin/uvid main.go

run: 
	GOFLAGS=-mod=mod go run main.go serve