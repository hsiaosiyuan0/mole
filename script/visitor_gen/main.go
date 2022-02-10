package main

import (
	"flag"
	"fmt"
	"log"
	"path"
	"path/filepath"
	"runtime"
)

func main() {
	_, fileName, _, _ := runtime.Caller(0)
	cwd, err := filepath.Abs(filepath.Dir(fileName))
	if err != nil {
		log.Fatal(err)
	}

	_ = path.Clean(path.Join(cwd, "..", "..", "ecma", "parser", "node_type.go"))
	// fmt.Println(cwd)

	// b, err := ioutil.ReadFile(file)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(b))
	fmt.Println(flag.Args())

	// fset := token.NewFileSet()
	// f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// d := f.Decls[0]
	// fmt.Println(d)

	// expr, err := parser.ParseExpr(`[visitor(before,after)]`)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(expr)
}
