package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
)

func TestLogicAssignFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_LOGIC_ASSIGN)
	TestFail(t, "a &&= b", "Unexpected token at (1:2)", opts)
}

func TestLogicAssignFail2(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_LOGIC_ASSIGN)
	TestFail(t, "a ||= b", "Unexpected token at (1:2)", opts)
}

func TestLogicAssignFail3(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_LOGIC_ASSIGN)
	TestFail(t, "a ??= b", "Unexpected token at (1:2)", opts)
}

func TestLogicAssignFail4(t *testing.T) {
	TestFail(t, "({a} &&= b)", "Assigning to rvalue at (1:1)", nil)
}

func TestLogicAssignFail5(t *testing.T) {
	TestFail(t, "({a} ||= b)", "Assigning to rvalue at (1:1)", nil)
}

func TestLogicAssignFail6(t *testing.T) {
	TestFail(t, "({a} ??= b)", "Assigning to rvalue at (1:1)", nil)
}
