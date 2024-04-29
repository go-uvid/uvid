test:
	cd js && pnpm --filter "uvid-js" test
	go test ./... -v --cover

test-report:
	go test ./... -v --cover -coverprofile=coverage.out
	go tool cover -html=coverage.out

build:
	GOFLAGS=-mod=mod go build

run:
	gin --immediate run main.go

publish-sdk:
	cd js/packages/uvid-js && npm version patch && pnpm publish
publish:
	goreleaser release --clean