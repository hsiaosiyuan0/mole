package analysis

import "github.com/hsiaosiyuan0/mole/ecma/parser"

type StubKind uint8

const (
	SK_NONE StubKind = 0
	SK_SEQ  StubKind = 1 << iota
	SK_JMP_FALSE
	SK_JMP_TRUE
	SK_LOOP
	SK_UNREACHABLE
)

type Block struct {
	Nodes []parser.Node
	In    []*Stub
	Out   []*Stub
}

type Stub struct {
	Kind StubKind
}
