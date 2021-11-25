package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/hsiaosiyuan0/mole/pkg/js/estree"
	"github.com/hsiaosiyuan0/mole/pkg/js/parser"
)

type Result struct {
	Data  interface{} `json:"data"`
	Error interface{} `json:"error"`
}

func compileToJson(src string) (string, error) {
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
	return string(b), nil
}

func compile(this js.Value, args []js.Value) interface{} {
	src := args[0].String()
	json, err := compileToJson(src)
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	return []interface{}{json, errMsg}
}

func main() {
	c := make(chan bool)
	js.Global().Set("compile", js.FuncOf(compile))
	<-c
}
