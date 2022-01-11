package parser

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/internal"
)

func TestTraverse(t *testing.T) {
	names := make([]string, 0)
	ast, err := compile("a + b - c", nil)
	AssertEqual(t, nil, err, "should be prog ok")

	ls := NewListenerImpl()
	ls[N_NAME_BEFORE] = func(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
		id := n.(*Ident)
		names = append(names, id.Text())
		return true
	}
	AssertEqual(t, true, DefaultListenerImpl[N_NAME_BEFORE] == nil, "should be copy")

	ListenProg(ast, ls, NewTraverseCtx())
	AssertEqual(t, 3, len(names), "should be 3 names")
}
