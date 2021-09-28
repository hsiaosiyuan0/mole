lint: cmd/lint/*.go internal/*/*.go pkg/*/*.go go.mod
	go build ./cmd/lint

install: ./lint
	cp ./lint /usr/local/bin/lint

clean:
	go run scripts/clean/mod.go

test: cmd/lint/*.go internal/*/*.go pkg/*/*.go go.mod
	go test ./...

bench-1: cmd/lint/*.go internal/*/*.go pkg/*/*.go go.mod
	go test -cpu 1 -benchmem -bench=. ./...
