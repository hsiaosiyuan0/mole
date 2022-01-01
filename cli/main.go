package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"

	"github.com/hsiaosiyuan0/mole/ecma/estree"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
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

	if *ast && strings.HasSuffix(*file, ".js") {
		printJsAst(string(src), *file)
	}
}

func assembleOutFile(file string, tag string) string {
	dir, f := path.Split(file)
	ext := path.Ext(f)
	fileNoExt := f[0:strings.LastIndex(f, ext)]
	nf := fileNoExt + tag + ext
	return path.Join(dir, nf)
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

	output := out.String()
	fmt.Println(output)

	outpath := assembleOutFile(file, "_ast")
	ioutil.WriteFile(outpath, out.Bytes(), 0644)
	fmt.Printf("Result is also saved in: %s\n", outpath)

	if ok := copyToClipboard(outpath); ok {
		fmt.Println("Result is also saved in clipboard")
	}

	return output, nil
}

func copyToClipboard(file string) bool {
	path, err := exec.LookPath("pbcopy")
	if err != nil {
		return false
	}

	cmd := exec.Command("bash", "-c", path+" < "+file)
	err = cmd.Run()
	return err == nil
}
