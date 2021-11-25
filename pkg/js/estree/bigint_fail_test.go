package estree

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func TestBigintFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "let i = 02n",
		"Identifier directly after number at (1:10)", opts)
}

func TestBigintFail2(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "i = 02n",
		"Identifier directly after number at (1:6)", opts)
}

func TestBigintFail3(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "((i = 02n) => {})",
		"Identifier directly after number at (1:8)", opts)
}

func TestBigintFail4(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "for (let i = 0n; i < 02n;++i) {}",
		"Identifier directly after number at (1:23)", opts)
}

func TestBigintFail5(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_STRICT)
	testFail(t, "i + 02n",
		"Identifier directly after number at (1:6)", opts)
}

func TestBigintFail6(t *testing.T) {
	testFail(t, "let i = 2e2n",
		"Identifier directly after number at (1:11)", nil)
}

func TestBigintFail7(t *testing.T) {
	testFail(t, "i = 2e2n",
		"Identifier directly after number at (1:7)", nil)
}

func TestBigintFail8(t *testing.T) {
	testFail(t, "((i = 2e2n) => {})",
		"Identifier directly after number at (1:9)", nil)
}

func TestBigintFail9(t *testing.T) {
	testFail(t, "for (let i = 0n; i < 2e2n;++i) {}",
		"Identifier directly after number at (1:24)", nil)
}

func TestBigintFail10(t *testing.T) {
	testFail(t, "i + 2e2n",
		"Identifier directly after number at (1:7)", nil)
}

func TestBigintFail11(t *testing.T) {
	testFail(t, "let i = 2.4n",
		"Identifier directly after number at (1:11)", nil)
}

func TestBigintFail12(t *testing.T) {
	testFail(t, "i = 2.4n",
		"Identifier directly after number at (1:7)", nil)
}

func TestBigintFail13(t *testing.T) {
	testFail(t, "((i = 2.4n) => {})",
		"Identifier directly after number at (1:9)", nil)
}

func TestBigintFail14(t *testing.T) {
	testFail(t, "for (let i = 0n; i < 2.4n;++i) {}",
		"Identifier directly after number at (1:24)", nil)
}

func TestBigintFail15(t *testing.T) {
	testFail(t, "i + 2.4n",
		"Identifier directly after number at (1:7)", nil)
}

func TestBigintFail16(t *testing.T) {
	testFail(t, "let i = .4n",
		"Identifier directly after number at (1:10)", nil)
}

func TestBigintFail17(t *testing.T) {
	testFail(t, "i = .4n",
		"Identifier directly after number at (1:6)", nil)
}

func TestBigintFail18(t *testing.T) {
	testFail(t, "((i = .4n) => {})",
		"Identifier directly after number at (1:8)", nil)
}

func TestBigintFail19(t *testing.T) {
	testFail(t, "for (let i = 0n; i < .4n;++i) {}",
		"Identifier directly after number at (1:23)", nil)
}

func TestBigintFail20(t *testing.T) {
	testFail(t, "i + .4n",
		"Identifier directly after number at (1:6)", nil)
}

func TestBigintFail21(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BIGINT)
	testFail(t, "let i = 0o2n",
		"Identifier directly after number at (1:11)", opts)
}

func TestBigintFail22(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BIGINT)
	testFail(t, "let i = 0b01n",
		"Identifier directly after number at (1:12)", opts)
}

func TestBigintFail23(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_BIGINT)
	testFail(t, "let i = 0x01n",
		"Identifier directly after number at (1:12)", opts)
}
