package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/hsiaosiyuan0/mole/pkg/js/estree"
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

func main() {
	file := flag.String("file", "", "the target file to be dealt with")
	ast := flag.Bool("ast", false, "print AST of the target file")
	flag.Parse()

	if *file == "" {
		panic("missing target file, use `-file` to specify one")
	}

	var src []byte
	var err error
	if src, err = ioutil.ReadFile(*file); err != nil {
		panic(err)
	}

	// js.TestLexer()
	if *ast && strings.HasSuffix(*file, ".js") {
		ast, err := printJsAst(string(src), *file)
		if err != nil {
			panic(err)
		}
		fmt.Println(ast)
	}
}

func printJsAst(src, file string) (string, error) {
	opts := parser.NewParserOpts()
	s := parser.NewSource("", src)
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
