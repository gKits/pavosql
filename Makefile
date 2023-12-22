build:
	@go build -o bin/pavosql cmd/pavosql/main.go

run:
	@go run cmd/pavosql/main.go

# Testing

test:
	@go test ./...

test.cover:
	@go test -coverprofile cover.out ./...

test.cover.show: test.cover
	@go tool cover -html cover.out

## Tree test

test.tree:
	@go test ./internal/tree/...

test.tree.cover:
	@go test -coverprofile tree_cover.out ./internal/tree/...

test.tree.cover.show: test.tree.cover
	@go tool cover -html tree_cover.out

## Parse test

test.parse:
	@go test ./internal/parse/...

test.parse.cover:
	@go test -coverprofile parse_cover.out ./internal/parse/...

test.parse.cover.show: test.parse.cover
	@go tool cover -html parse_cover.out

# Cleanup

cleancover:
	@rm *cover.out
