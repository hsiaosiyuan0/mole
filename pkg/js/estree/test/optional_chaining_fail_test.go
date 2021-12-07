package estree_test

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func TestOptionalChainFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_OPT_EXPR)
	testFail(t, "obj?.foo", "Unexpected token at (1:3)", opts)
}

func TestOptionalChainFail2(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_OPT_EXPR)
	testFail(t, "obj?.[foo]", "Unexpected token at (1:3)", opts)
}

func TestOptionalChainFail3(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_OPT_EXPR)
	testFail(t, "obj?.()", "Unexpected token at (1:3)", opts)
}

func TestOptionalChainFail4(t *testing.T) {
	testFail(t, "obj?.0", "Unexpected token `EOF` at (1:6)", nil)
}

func TestOptionalChainFail5(t *testing.T) {
	testFail(t, "async?.() => {}", "Unexpected token at (1:10)", nil)
}

func TestOptionalChainFail6(t *testing.T) {
	testFail(t, "new obj?.()",
		"Invalid optional chain from new expression at (1:7)", nil)
}

func TestOptionalChainFail7(t *testing.T) {
	testFail(t, "new obj?.foo()",
		"Invalid optional chain from new expression at (1:7)", nil)
}

func TestOptionalChainFail8(t *testing.T) {
	testFail(t, "obj?.foo\n`template`",
		"Invalid tagged template on optional chain at (2:0)", nil)
}

func TestOptionalChainFail9(t *testing.T) {
	testFail(t, "obj?.foo = 0",
		"Assigning to rvalue at (1:0)", nil)
}

func TestOptionalChainFail10(t *testing.T) {
	testFail(t, "obj?.foo.bar = 0",
		"Assigning to rvalue at (1:0)", nil)
}

func TestOptionalChainFail11(t *testing.T) {
	testFail(t, "obj?.().foo = 0",
		"Assigning to rvalue at (1:0)", nil)
}

func TestOptionalChainFail12(t *testing.T) {
	testFail(t, "obj?.foo++",
		"Assigning to rvalue at (1:0)", nil)
}

func TestOptionalChainFail13(t *testing.T) {
	testFail(t, "obj?.foo--",
		"Assigning to rvalue at (1:0)", nil)
}

func TestOptionalChainFail14(t *testing.T) {
	testFail(t, "++obj?.foo",
		"Assigning to rvalue at (1:2)", nil)
}

func TestOptionalChainFail15(t *testing.T) {
	testFail(t, "--obj?.foo",
		"Assigning to rvalue at (1:2)", nil)
}

func TestOptionalChainFail16(t *testing.T) {
	testFail(t, "obj?.foo.bar++",
		"Assigning to rvalue at (1:0)", nil)
}

func TestOptionalChainFail17(t *testing.T) {
	testFail(t, "for (obj?.foo in {});",
		"Assigning to rvalue at (1:5)", nil)
}

func TestOptionalChainFail18(t *testing.T) {
	testFail(t, "for (obj?.foo.bar in {});",
		"Assigning to rvalue at (1:5)", nil)
}

func TestOptionalChainFail19(t *testing.T) {
	testFail(t, "for (obj?.foo of []);",
		"Assigning to rvalue at (1:5)", nil)
}

func TestOptionalChainFail20(t *testing.T) {
	testFail(t, "for (obj?.foo.bar of []);",
		"Assigning to rvalue at (1:5)", nil)
}
