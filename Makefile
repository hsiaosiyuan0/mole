lint: cmd/lint/*.go internal/*/*.go pkg/*/*.go go.mod
	go build ./cmd/lint

install: ./lint
	cp ./lint /usr/local/bin/lint

clean:
	go run scripts/clean/mod.go

test: cmd/lint/*.go internal/*/*.go go.mod pkg/*/*.go
	go test ./...
