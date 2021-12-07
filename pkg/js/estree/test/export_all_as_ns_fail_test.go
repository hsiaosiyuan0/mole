package estree_test

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func TestExportAllAsNSFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_EXPORT_ALL_AS_NS)
	testFail(t, "export * as ns from \"source\"", "Unexpected token `identifier` at (1:9)", opts)
}

func TestExportAllAsNSFail2(t *testing.T) {
	testFail(t, "export * as ns", "Unexpected token `EOF` at (1:14)", nil)
}

func TestExportAllAsNSFail3(t *testing.T) {
	testFail(t, "export * as from \"source\"", "Unexpected token `string` at (1:17)", nil)
}

func TestExportAllAsNSFail4(t *testing.T) {
	testFail(t, "export * as ns \"source\"", "Unexpected token `string` at (1:15)", nil)
}

func TestExportAllAsNSFail5(t *testing.T) {
	testFail(t, "export {} as ns from \"source\"", "Unexpected token at (1:10)", nil)
}
