package estree_test

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func TestNumSepFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_NUM_SEP)
	testFail(t, "123_456", "Invalid number at (1:3)", opts)
}

func TestNumSepFail2(t *testing.T) {
	testFail(t, "123__456",
		"Only one underscore is allowed as numeric separator at (1:4)", nil)
}

func TestNumSepFail3(t *testing.T) {
	testFail(t, "0._123456",
		"Numeric separator is not allowed at the first of digits at (1:2)", nil)
}

func TestNumSepFail4(t *testing.T) {
	testFail(t, "123456_",
		"Numeric separator is not allowed at the last of digits at (1:6)", nil)
}

func TestNumSepFail5(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "012_345",
		"Numeric separator is not allowed in legacy octal numeric literals at (1:3)", opts)
}

func TestNumSepFail6(t *testing.T) {
	testFail(t, "'\\x2_0'",
		"Bad character escape sequence at (1:2)", nil)
}

func TestNumSepFail7(t *testing.T) {
	testFail(t, "'\\u00_20'", "Bad character escape sequence at (1:2)", nil)
}

func TestNumSepFail8(t *testing.T) {
	testFail(t, "'\\u{2_0}'", "Bad character escape sequence at (1:2)", nil)
}

func TestNumSepFail9(t *testing.T) {
	testFail(t, "0b_10",
		"Numeric separator is not allowed at the first of digits at (1:2)", nil)
}

func TestNumSepFail10(t *testing.T) {
	testFail(t, "0b10_",
		"Numeric separator is not allowed at the last of digits at (1:4)", nil)
}

func TestNumSepFail11(t *testing.T) {
	testFail(t, "0b10__10",
		"Only one underscore is allowed as numeric separator at (1:5)", nil)
}

func TestNumSepFail12(t *testing.T) {
	testFail(t, "0o_7",
		"Numeric separator is not allowed at the first of digits at (1:2)", nil)
}

func TestNumSepFail13(t *testing.T) {
	testFail(t, "0o7_",
		"Numeric separator is not allowed at the last of digits at (1:3)", nil)
}

func TestNumSepFail14(t *testing.T) {
	testFail(t, "0o7__07",
		"Only one underscore is allowed as numeric separator at (1:4)", nil)
}

func TestNumSepFail15(t *testing.T) {
	testFail(t, "0x_a",
		"Numeric separator is not allowed at the first of digits at (1:2)", nil)
}

func TestNumSepFail16(t *testing.T) {
	testFail(t, "0xa_",
		"Numeric separator is not allowed at the first of digits at (1:3)", nil)
}

func TestNumSepFail17(t *testing.T) {
	testFail(t, "0xa__a",
		"Only one underscore is allowed as numeric separator at (1:4)", nil)
}
