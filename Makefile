test:
	pnpm --filter "./packages/**" test
	go test ./... -v --cover

test-report:
	go test ./... -v --cover -coverprofile=coverage.out
	go tool cover -html=coverage.out

build:
	pnpm --filter "./packages/**" build
	GOFLAGS=-mod=mod go build

run:
	gin --immediate run main.go

publish-sdk:
	pnpm --filter "uvid-js" exec npm version patch && pnpm --filter "uvid-js" publish
	
publish:
	goreleaser release --clean