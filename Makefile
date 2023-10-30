build:
	@go build -o bin/pavosql ./cmd/pavosql

run: build
	@./bin/pavosql

test:
	@go test -v ./...
