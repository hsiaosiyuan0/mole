package estree_test

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func TestImportMetaFail1(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_META_PROPERTY)
	testFail(t, "import.meta", "Unexpected token `.` at (1:6)", opts)
}

func TestImportMetaFail2(t *testing.T) {
	// it's not a requirement in mole to set FEAT_MODULE to enable the `import.meta` syntax
	// use `FEAT_META_PROPERTY` to turn on/off the syntax instead
	// testFail(t, "import.meta", "Cannot use 'import.meta' outside a module (1:0)", nil)
}

func TestImportMetaFail3(t *testing.T) {
	testFail(t, "import['meta']", "Unexpected token `[` at (1:6)", nil)
}

func TestImportMetaFail4(t *testing.T) {
	testFail(t, "a = import['meta']", "Unexpected token `[` at (1:10)", nil)
}

func TestImportMetaFail5(t *testing.T) {
	testFail(t, "import.target",
		"The only valid meta property for import is `import.meta` at (1:7)", nil)
}

func TestImportMetaFail6(t *testing.T) {
	testFail(t, "new.meta",
		"The only valid meta property for new is `new.target` at (1:4)", nil)
}

func TestImportMetaFail7(t *testing.T) {
	testFail(t, "im\\u0070ort.meta",
		"Keyword must not contain escaped characters at (1:0)", nil)
}

func TestImportMetaFail8(t *testing.T) {
	testFail(t, "import.\\u006d\\u0065\\u0074\\u0061",
		"Meta property can not contain escaped characters at (1:7)", nil)
}
