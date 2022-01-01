package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
)

func TestNullishFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_NULLISH)
	TestFail(t, "a ?? b", "Unexpected token at (1:2)", opts)
}

func TestNullishFail2(t *testing.T) {
	TestFail(t, "?? b", "Unexpected token `??` at (1:0)", nil)
}

func TestNullishFail3(t *testing.T) {
	TestFail(t, "a ??", "Unexpected token `EOF` at (1:4)", nil)
}

func TestNullishFail4(t *testing.T) {
	TestFail(t, "a || b ?? c",
		"Cannot use unparenthesized `??` within logic expressions at (1:7)", nil)
}

func TestNullishFail5(t *testing.T) {
	TestFail(t, "a && b ?? c",
		"Cannot use unparenthesized `??` within logic expressions at (1:7)", nil)
}

func TestNullishFail6(t *testing.T) {
	TestFail(t, "a ?? b || c",
		"Cannot use unparenthesized `??` within logic expressions at (1:7)", nil)
}

func TestNullishFail7(t *testing.T) {
	TestFail(t, "a ?? b && c",
		"Cannot use unparenthesized `??` within logic expressions at (1:7)", nil)
}

func TestNullishFail8(t *testing.T) {
	TestFail(t, "a+1 || b+1 ?? c",
		"Cannot use unparenthesized `??` within logic expressions at (1:11)", nil)
}

func TestNullishFail9(t *testing.T) {
	TestFail(t, "a+1 && b+1 ?? c",
		"Cannot use unparenthesized `??` within logic expressions at (1:11)", nil)
}

func TestNullishFail10(t *testing.T) {
	TestFail(t, "a+1 ?? b+1 || c",
		"Cannot use unparenthesized `??` within logic expressions at (1:11)", nil)
}

func TestNullishFail11(t *testing.T) {
	TestFail(t, "a+1 ?? b+1 && c",
		"Cannot use unparenthesized `??` within logic expressions at (1:11)", nil)
}
