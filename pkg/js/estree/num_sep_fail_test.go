package estree

import "testing"

func TestNumSepFail1(t *testing.T) {
	testFail(t, "class A { async* f() { () => await a; } }",
		"Cannot use keyword 'await' outside an async function at (1:29)", nil)
}

func TestNumSepFail2(t *testing.T) {
	testFail(t, "class A { async* f() { () => await a; } }",
		"Cannot use keyword 'await' outside an async function at (1:29)", nil)
}

func TestNumSepFail3(t *testing.T) {
	testFail(t, "class A { async* f() { () => await a; } }",
		"Cannot use keyword 'await' outside an async function at (1:29)", nil)
}

func TestNumSepFail4(t *testing.T) {
	testFail(t, "class A { async* f() { () => await a; } }",
		"Cannot use keyword 'await' outside an async function at (1:29)", nil)
}

func TestNumSepFail5(t *testing.T) {
	testFail(t, "class A { async* f() { () => await a; } }",
		"Cannot use keyword 'await' outside an async function at (1:29)", nil)
}

func TestNumSepFail6(t *testing.T) {
	testFail(t, "class A { async* f() { () => await a; } }",
		"Cannot use keyword 'await' outside an async function at (1:29)", nil)
}

func TestNumSepFail7(t *testing.T) {
	testFail(t, "class A { async* f() { () => await a; } }",
		"Cannot use keyword 'await' outside an async function at (1:29)", nil)
}

func TestNumSepFail8(t *testing.T) {
	testFail(t, "class A { async* f() { () => await a; } }",
		"Cannot use keyword 'await' outside an async function at (1:29)", nil)
}
