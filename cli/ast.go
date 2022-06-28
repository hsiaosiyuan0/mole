package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hsiaosiyuan0/mole/ecma/estree"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/span"
)

type AstInspector struct {
}

func (a *AstInspector) Process(opts *Options) bool {
	if !opts.ast {
		return false
	}

	if opts.file == "" {
		panic("missing target file, use `-file` to specify one")
	}

	var src []byte
	var err error
	if src, err = ioutil.ReadFile(opts.file); err != nil {
		panic(err)
	}

	if opts.ast && strings.HasSuffix(opts.file, ".js") {
		printJsAst(string(src), opts.file)
	}

	return true
}

func assembleOutFile(file string, tag string) string {
	dir, f := filepath.Split(file)
	ext := filepath.Ext(f)
	fileNoExt := f[0:strings.LastIndex(f, ext)]
	nf := fileNoExt + tag + ext
	return filepath.Join(dir, nf)
}

func printJsAst(src, file string) (string, error) {
	opts := parser.NewParserOpts()
	s := span.NewSource("", src)
	p := parser.NewParser(s, opts)
	ast, err := p.Prog()
	if err != nil {
		return "", err
	}

	b, err := json.Marshal(estree.ConvertProg(ast.(*parser.Prog), estree.NewConvertCtx()))
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")

	output := out.String()
	fmt.Println(output)

	outFile := assembleOutFile(file, "_ast")
	ioutil.WriteFile(outFile, out.Bytes(), 0644)
	fmt.Printf("Result is also saved in: %s\n", outFile)

	if ok := copyToClipboard(outFile); ok {
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
