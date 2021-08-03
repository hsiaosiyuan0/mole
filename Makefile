mlint: cmd/mlint/*.go internal/*/*.go go.mod
	go build ./cmd/mlint

install: ./mlint
	cp ./mlint /usr/local/bin/mlint

clean:
	go run scripts/clean/mod.go

test: cmd/mlint/*.go internal/*/*.go go.mod
	go test ./...
