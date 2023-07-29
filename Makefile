test:
	cd js && pnpm --filter "uvid-js" test
	go test ./... -v --cover

test-report:
	go test ./... -v --cover -coverprofile=coverage.out
	go tool cover -html=coverage.out

build:
	cd js && pnpm --filter "uvid-js" --filter "dash" build
	GOFLAGS=-mod=mod go build -o bin/uvid main.go

run: 
	gin --immediate run main.go

publish-sdk:
	cd js/packages/uvid-js && pnpm version patch && pnpm publish
publish:
	goreleaser release