package parser

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func TestTraverse(t *testing.T) {
	names := make([]string, 0)
	ast, err := compile("a + b - c", nil)
	assert.Equal(t, nil, err, "should be prog ok")

	ls := NewListenerImpl()
	ls[N_NAME_BEFORE] = func(n Node, v interface{}, ctx *TraverseCtx) ContOrStop {
		id := n.(*Ident)
		names = append(names, id.Text())
		return true
	}
	assert.Equal(t, true, DefaultListenerImpl[N_NAME_BEFORE] == nil, "should be copy")

	ListenProg(ast, ls, NewTraverseCtx())
	assert.Equal(t, 3, len(names), "should be 3 names")
}
