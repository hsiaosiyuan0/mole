package estree_test

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func TestNullishFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_NULLISH)
	testFail(t, "a ?? b", "Unexpected token at (1:2)", opts)
}

func TestNullishFail2(t *testing.T) {
	testFail(t, "?? b", "Unexpected token `??` at (1:0)", nil)
}

func TestNullishFail3(t *testing.T) {
	testFail(t, "a ??", "Unexpected token `EOF` at (1:4)", nil)
}

func TestNullishFail4(t *testing.T) {
	testFail(t, "a || b ?? c",
		"Cannot use unparenthesized `??` within logic expressions at (1:7)", nil)
}

func TestNullishFail5(t *testing.T) {
	testFail(t, "a && b ?? c",
		"Cannot use unparenthesized `??` within logic expressions at (1:7)", nil)
}

func TestNullishFail6(t *testing.T) {
	testFail(t, "a ?? b || c",
		"Cannot use unparenthesized `??` within logic expressions at (1:7)", nil)
}

func TestNullishFail7(t *testing.T) {
	testFail(t, "a ?? b && c",
		"Cannot use unparenthesized `??` within logic expressions at (1:7)", nil)
}

func TestNullishFail8(t *testing.T) {
	testFail(t, "a+1 || b+1 ?? c",
		"Cannot use unparenthesized `??` within logic expressions at (1:11)", nil)
}

func TestNullishFail9(t *testing.T) {
	testFail(t, "a+1 && b+1 ?? c",
		"Cannot use unparenthesized `??` within logic expressions at (1:11)", nil)
}

func TestNullishFail10(t *testing.T) {
	testFail(t, "a+1 ?? b+1 || c",
		"Cannot use unparenthesized `??` within logic expressions at (1:11)", nil)
}

func TestNullishFail11(t *testing.T) {
	testFail(t, "a+1 ?? b+1 && c",
		"Cannot use unparenthesized `??` within logic expressions at (1:11)", nil)
}
