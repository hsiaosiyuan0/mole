package parser

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
)

func TestJSX(t *testing.T) {
	ast, err := compile(`
  <div>test</div>
  `, nil)
	assert.Equal(t, nil, err, "should be prog ok")

	fn := ast.(*Prog).stmts[0].(*FnDec)
	id := fn.id.(*Ident)
	assert.Equal(t, "a", id.Text(), "should be a")
}
