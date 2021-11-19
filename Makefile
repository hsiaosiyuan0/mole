dep: cmd/mole/*.go pkg/*/*.go go.mod
	go get -v ./...

mole: cmd/mole/*.go pkg/*/*.go go.mod
	go build ./cmd/mole

install: ./mole
	cp ./mole /usr/local/bin/mole

clean:
	go run scripts/clean/mod.go

test: cmd/mole/*.go pkg/*/*.go go.mod
	go test ./...

test-estree-convert: cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree  -run "^Test\d"

test-harmony:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree -run "^TestHarmony\d"

test-harmony-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree -run "^TestHarmonyFail\d"

test-async-iteration:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree -run "^TestAsyncIteration\d"

test-async-iteration-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree -run "^TestAsyncIterationFail\d"

test-async-await:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree -run "^TestAsyncAwait\d"

test-async-await-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree -run "^TestAsyncAwaitFail\d"

test-parser: cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/parser

bench-1: cmd/mole/*.go pkg/*/*.go go.mod
	go test -cpu 1 -benchmem -bench=. ./...
