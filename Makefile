dep: cmd/mole/*.go pkg/*/*.go go.mod
	go get -v ./...

mole: cmd/mole/*.go pkg/*/*.go go.mod
	go build ./cmd/mole

mole_wasm: cmd/mole_wasm/*.go pkg/*/*.go go.mod
	GOOS=js GOARCH=wasm go build -o mole.wasm ./cmd/mole_wasm

install: ./mole
	cp ./mole /usr/local/bin/mole

clean:
	go run scripts/clean/mod.go

test: cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/...

test-estree-convert: cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/...  -run "^Test\d"

test-harmony:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestHarmony\d"

test-harmony-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestHarmonyFail\d"

test-async-iteration:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestAsyncIteration\d"

test-async-iteration-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestAsyncIterationFail\d"

test-async-await:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estre/...e -run "^TestAsyncAwait\d"

test-async-await-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestAsyncAwaitFail\d"

test-class:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestClassFeature\d"

test-class-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestClassFeatureFail\d"

test-optional-chain:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestOptionalChain\d"

test-optional-chain-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestOptionalChainFail\d"

test-es7:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestEs7th\d"

test-es7-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestEs7thFail\d"

test-nullish:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestNullish\d"

test-nullish-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestNullishFail\d"

test-directive:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestDirective\d"

test-num-sep:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestNumSep\d"

test-num-sep-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestNumSepFail\d"

test-logic-assign:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestLogicAssign\d"

test-logic-assign-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestLogicAssignFail\d"

test-dynamic-import:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestDynamicImport\d"

test-dynamic-import-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestDynamicImportFail\d"

test-bigint:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestBigint\d"

test-bigint-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestBigintFail\d"

test-jsx:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestJSX\d"

test-jsx-fail:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestJSXFail\d"

test-export-all-as-ns:	cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestExportAllAsNS"

test-fixtures: cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/estree/... -run "^TestFixtures_"

test-parser: cmd/mole/*.go pkg/*/*.go go.mod
	go test ./pkg/js/parser

bench-1: cmd/mole/*.go pkg/*/*.go go.mod
	go test -cpu 1 -benchmem -bench=. ./pkg/js/estree/... -run "^Benchmark"

html-entities:
	go run scripts/html_entities/mod.go

gofmt:
	gofmt -w .