package estree

import (
	"testing"

	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func TestDynamicImportFail1(t *testing.T) {
	testFail(t, "function failsParse() { return import.then(); }",
		"The only valid meta property for import is `import.meta` at (1:38)", nil)
}

func TestDynamicImportFail2(t *testing.T) {
	testFail(t, "var dynImport = import; dynImport('http');",
		"Unexpected token `;` at (1:22)", nil)
}

func TestDynamicImportFail3(t *testing.T) {
	opts := parser.NewParserOpts()
	opts.Feature = opts.Feature.Off(parser.FEAT_DYNAMIC_IMPORT)
	testFail(t, "import('test.js')", "Unexpected token `(` at (1:6)", opts)
}

func TestDynamicImportFail4(t *testing.T) {
	testFail(t, "import()", "Unexpected token `)` at (1:7)", nil)
}

func TestDynamicImportFail5(t *testing.T) {
	testFail(t, "import(a, b)", "Unexpected token `,` at (1:8)", nil)
}

func TestDynamicImportFail6(t *testing.T) {
	testFail(t, "import(...[a])", "Unexpected token `...` at (1:7)", nil)
}

func TestDynamicImportFail7(t *testing.T) {
	testFail(t, "import(source,)", "Unexpected token `,` at (1:13)", nil)
}

func TestDynamicImportFail8(t *testing.T) {
	testFail(t, "new import(source)", "Cannot use new with `import()` at (1:4)", nil)
}

func TestDynamicImportFail9(t *testing.T) {
	testFail(t, "(import)(s)", "Unexpected token `)` at (1:7)", nil)
}
