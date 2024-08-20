build: 
	@go build -o bin/api

run: build
	@./bin/api

dev: 
	air


test:
	@go test -v ./...

seed:
	@go run scripts/seed.go