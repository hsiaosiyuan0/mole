package estree_test

import (
	"testing"

	. "github.com/hsiaosiyuan0/mole/ecma/estree/test"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
)

func TestFixture_core(t *testing.T) {
	RunFixtures(t, "core", parser.NewParserOpts())
}

func TestFixture_es2015(t *testing.T) {
	RunFixtures(t, "es2015", parser.NewParserOpts())
}

func TestFixture_ts(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.On(parser.FEAT_TS)
	// `FEAT_JSX` will be turned on by the `options.json` of the
	// test case automatically
	opts.Feature = opts.Feature.Off(parser.FEAT_JSX)
	RunFixtures(t, "typescript", opts)
}

func TestFixture_tsManually(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.On(parser.FEAT_TS)
	// `FEAT_JSX` will be turned on by the `options.json` of the
	// test case automatically
	opts.Feature = opts.Feature.Off(parser.FEAT_JSX)
	RunFixtures(t, "typescript/type-arguments-bit-shift-left-like/after-bit-shift", opts)
}
