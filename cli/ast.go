package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hsiaosiyuan0/mole/ecma/estree"
	"github.com/hsiaosiyuan0/mole/ecma/parser"
	"github.com/hsiaosiyuan0/mole/span"

	_runtime "runtime/pprof"
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

	if opts.ast {
		printJsAst(string(src), opts.file, opts.perf)
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

func printJsAst(src, file string, perf bool) (string, error) {
	opts := parser.NewParserOpts()
	s := span.NewSource("", src)
	p := parser.NewParser(s, opts)

	var ast parser.Node
	var err error
	if perf {
		cpuProf, err := os.Create("cpu.prof")
		if err != nil {
			return "", err
		}
		defer cpuProf.Close()

		heap, err := os.Create("heap.out")
		if err != nil {
			return "", err
		}
		defer cpuProf.Close()

		_runtime.StartCPUProfile(cpuProf)
		ast, err = p.Prog()
		_runtime.StopCPUProfile()
		_runtime.WriteHeapProfile(heap)
	} else {
		start := time.Now()
		ast, err = p.Prog()
		elapsed := time.Since(start)
		fmt.Printf("Parsed in %dμs\n", elapsed.Microseconds())
	}

	if err != nil {
		return "", err
	}

	ctx := estree.NewConvertCtx(p)
	ctx.LineCol = false
	b, err := json.Marshal(estree.ConvertProg(ast.(*parser.Prog), ctx))
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")

	output := out.String()

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
