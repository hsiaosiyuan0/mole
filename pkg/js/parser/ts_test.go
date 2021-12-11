package parser

import (
	"testing"
)

func compileTs(code string, opts *ParserOpts) (Node, error) {
	if opts == nil {
		opts = NewParserOpts()
	}
	opts.Feature = opts.Feature.On(FEAT_TS)
	p := newParser(code, opts)
	return p.Prog()
}

func TestTs(t *testing.T) {
	// ast, err := compileTs("var a: (string|number) = 1", nil)
	// assert.Equal(t, nil, err, "should be prog ok")
	// fmt.Println(ast)
}
