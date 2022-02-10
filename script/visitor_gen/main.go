package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
)

func main() {
	_, fileName, _, _ := runtime.Caller(0)
	cwd, err := filepath.Abs(filepath.Dir(fileName))
	if err != nil {
		log.Fatal(err)
	}

	file := path.Clean(path.Join(cwd, "..", "..", "ecma", "parser", "node_type.go"))
	fmt.Println(cwd)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)

	if err != nil {
		log.Fatal(err)
	}

	var bf bytes.Buffer
	ast.Fprint(&bf, fset, f, func(string, reflect.Value) bool {
		return true
	})

	fmt.Println(bf.String())

	d := f.Decls[0]
	fmt.Println(d)

	expr, err := parser.ParseExpr(`[visitor(before,after)]`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(expr)
}
