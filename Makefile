dep:
	go get -v ./...

mole:
	go build -o mole ./cli

mole_wasm:
	GOOS=js GOARCH=wasm go build -o mole.wasm ./wasm

install:
	cp ./mole /usr/local/bin/mole

clean:
	go clean -cache
	rm -rf mole
	rm -rf mole.wasm

test:
	go clean -testcache
	go test ./ecma/...

test-ecma:
	go test ./ecma/...

test-estree-basic:
	go test ./ecma/estree/test/basic... -run "^Test"

test-estree-fixture:
	go test ./ecma/estree/test/fixture... -run "^TestFixture"

bench-ecma:
	go test -cpu 1 -benchmem -bench=. ./ecma/estree/test/perf... -run "^Benchmark"

html-entities:
	go run script/html_entities/main.go

gofmt:
	gofmt -w .