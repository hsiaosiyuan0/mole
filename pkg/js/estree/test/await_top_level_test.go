package estree_test

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/assert"
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func TestAwaitTopLevel1(t *testing.T) {
	ast, err := compile("await 1")
	assert.Equal(t, nil, err, "should be prog ok")

	assert.EqualJson(t, `
{
  "type": "Program",
  "start": 0,
  "end": 7,
  "body": [
    {
      "type": "ExpressionStatement",
      "start": 0,
      "end": 7,
      "expression": {
        "type": "AwaitExpression",
        "start": 0,
        "end": 7,
        "argument": {
          "type": "Literal",
          "start": 6,
          "end": 7,
          "value": 1
        }
      }
    }
  ]
}
`, ast)
}

func TestAwaitTopLevelFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_MODULE).Off(parser.FEAT_GLOBAL_ASYNC)
	testFail(t, "await 1", "Unexpected token at (1:6)", opts)
}

func TestAwaitTopLevelFail2(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_MODULE).Off(parser.FEAT_GLOBAL_ASYNC)
	testFail(t, "function foo() {return await 1}", "Unexpected token at (1:29)", opts)
}

func TestAwaitTopLevelFail3(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_GLOBAL_ASYNC)
	testFail(t, "function foo() {return await 1}",
		"Cannot use keyword 'await' outside an async function at (1:23)", nil)
}

func TestAwaitTopLevelFail4(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_GLOBAL_ASYNC)
	testFail(t, "await 1", "Cannot use keyword 'await' outside an async function at (1:0)", opts)
}
