package estree_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/hsiaosiyuan0/mole/ecma/estree"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/internal"
)

func NewParser(code string, opts *parser.ParserOpts) *parser.Parser {
	if opts == nil {
		opts = parser.NewParserOpts()
	}
	s := parser.NewSource("", code)
	return parser.NewParser(s, opts)
}

func CompileProg(code string, opts *parser.ParserOpts) (parser.Node, error) {
	p := NewParser(code, opts)
	return p.Prog()
}

func Compile(code string) (string, error) {
	return CompileWithOpts(code, parser.NewParserOpts())
}

func CompileWithOpts(code string, opts *parser.ParserOpts) (string, error) {
	s := parser.NewSource("", code)
	p := parser.NewParser(s, opts)
	ast, err := p.Prog()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(estree.ConvertProg(ast.(*parser.Prog)))
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")

	return out.String(), nil
}

func TestFail(t *testing.T, code, errMs string, opts *parser.ParserOpts) {
	ast, err := CompileProg(code, opts)
	if err == nil {
		t.Fatalf("should not pass code:\n%s\nast:\n%v", code, ast)
	}
	internal.AssertEqual(t, errMs, err.Error(), "")
}
