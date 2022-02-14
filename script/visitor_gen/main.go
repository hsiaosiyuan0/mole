package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/hsiaosiyuan0/mole/fuzz"
	"github.com/hsiaosiyuan0/mole/script/macro"
)

type MethodInfo struct {
	Name  string
	Nodes bool // whether the getter returns `[]parser.Node` or not
	Dec   *ast.FuncType
}

type TplParamsGenVisitor struct {
	Name    string
	TypName string
	Methods []*MethodInfo
}

func genVisitor(output io.Writer, structName string, typName string, m []*MethodInfo) error {
	name := fmt.Sprintf("visitor_%s", structName)
	// use tpl instead of manually construct the ast
	fnMap := template.FuncMap{
		"UnPrefix": unPrefix,
	}
	tpl, err := template.New(name).Funcs(fnMap).Parse(`
func Visit{{ .Name }}(node parser.Node, ctx *WalkCtx) {
  {{- if ge (len .Methods) 1 }}
    n := node.(*parser.{{ .Name }})

    CallVisitor(VK_{{ .TypName | UnPrefix }}_BEFORE, n, ctx)
    {{- range $key, $value := .Methods }}
      {{ if .Nodes }}
        VisitNodes(n.{{ $value.Name }}(), ctx)
      {{- else }}
        VisitNode(n.{{ $value.Name }}(), ctx)
      {{- end }}
      if ctx.stop {
        return
      }
    {{ end }}
    CallVisitor(VK_{{ .TypName | UnPrefix }}_AFTER, n, ctx)
  {{ end }}
}
  `)
	if err != nil {
		return err
	}
	return tpl.Execute(output, &TplParamsGenVisitor{structName, typName, m})
}

func genVisitors(output io.Writer, nodeTypStruct map[string]string, structMethods map[string][]*MethodInfo) {
	output.Write([]byte(`
type Visitor = func(node parser.Node, ctx *WalkCtx)
type Visitors = [VK_DEF_END][]Visitor

func AddVisitor(vs *Visitors, vk VisitorKind, impl Visitor) {
  hs := vs[vk]
  if hs == nil {
    hs = []Visitor{}
    vs[vk] = hs
  }
  vs[vk] = append(hs, impl)
}
  `))

	processedVisitors := map[string]bool{}
	for nodeTyp, structName := range nodeTypStruct {
		if _, ok := processedVisitors[structName]; ok {
			continue
		}
		err := genVisitor(output, structName, nodeTyp, structMethods[structName])
		if err != nil {
			log.Fatal(err)
		}
		processedVisitors[structName] = true
	}
}

func unPrefix(s string) string {
	idx := strings.IndexByte(s, '_')
	return s[idx+1:]
}

type TplParamsGenVisitorKinds struct {
	NodeTypStruct map[string]string
	StructMethods map[string][]*MethodInfo
}

func genVisitorKinds(output io.Writer, nodeTypStruct map[string]string, structMethods map[string][]*MethodInfo) error {
	fnMap := template.FuncMap{
		"ToUpper":  strings.ToUpper,
		"UnPrefix": unPrefix,
	}
	tpl, err := template.New("visitor types").Funcs(fnMap).Parse(`
type VisitorKind uint16

const (
  VK_ILLEGAL VisitorKind = 0

  {{- range $key, $value := .NodeTypStruct }}
    VK_{{- $key | UnPrefix | ToUpper }} = VisitorKind(parser.{{ $key }})
  {{- end }}

  VK_BEFORE_AFTER = VisitorKind(parser.N_NODE_DEF_END) + iota
  {{- range $key, $value := .NodeTypStruct }}
    {{- if ge (len (index $.StructMethods $value)) 1  }}
      VK_{{ $key | UnPrefix | ToUpper }}_BEFORE
      VK_{{ $key | UnPrefix | ToUpper }}_AFTER
    {{- end }}
  {{- end }}

  VK_DEF_END
)
  `)
	if err != nil {
		return err
	}
	return tpl.Execute(output, &TplParamsGenVisitorKinds{nodeTypStruct, structMethods})
}

func genDefaultVisitors(output io.Writer, nodeTypStruct map[string]string) error {
	fnMap := template.FuncMap{
		"ToUpper":  strings.ToUpper,
		"UnPrefix": unPrefix,
	}
	tpl, err := template.New("visitor types").Funcs(fnMap).Parse(`
var DefaultVisitors Visitors = [VK_DEF_END][]Visitor{}

func init() {
  {{- range $key, $value := . }}
    DefaultVisitors[VK_{{ $key | UnPrefix | ToUpper }}] = []Visitor{Visit{{ $value }}}
  {{- end }}
}
  `)
	if err != nil {
		return err
	}
	return tpl.Execute(output, nodeTypStruct)
}

func IfNodesReurned(f *ast.FuncType, name string) bool {
	if f.Results == nil || len(f.Results.List) != 1 {
		log.Fatalf("%s should return only one value", name)
	}
	r := f.Results.List[0]
	a, ok := r.Type.(*ast.ArrayType)
	if !ok {
		return false
	}
	i, ok := a.Elt.(*ast.Ident)
	if !ok {
		return false
	}
	return i.Name == "Node"
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	distFile := path.Join(wd, "visitor.go")
	_, err = os.Stat(distFile)
	if err != nil {
		return
	}

	var defDir string
	flag.StringVar(&defDir, "d", "", "the AST definition directory, relative with current file")
	flag.Parse()

	ctxs, procCtx, err := macro.MacroCtxsOfWorkingDir(wd, defDir)
	if err != nil {
		log.Fatal(err)
	}

	nodeTypStruct := map[string]string{}
	structMethods := map[string][]*MethodInfo{}
	for _, ctx := range ctxs {
		if v, ok := ctx.Node.(*ast.ValueSpec); ok {
			nodeTyp := v.Names[0].Name
			structName := ctx.Args[0].(string)
			nodeTypStruct[nodeTyp] = structName
		} else if name, _, ok := macro.IsStructDec(ctx.Node); ok {
			ns := []*MethodInfo{}
			for _, n := range ctx.Args {
				ns = append(ns, &MethodInfo{Name: n.(string), Dec: nil})
			}
			structMethods[name] = ns
		}
	}

	// walk the pkgs to find out the type of `MethodInfo.dec`
	macro.WalkPkgs(procCtx.Pkgs, func(f *ast.File, s string, pc *macro.ProcCtx) error {
		for _, dec := range f.Decls {
			if v, ok := dec.(*ast.FuncDecl); ok {
				if v.Recv == nil {
					continue
				}
				name := v.Name.Name
				fnTyp := v.Type
				if v, ok := v.Recv.List[0].Type.(*ast.StarExpr); ok {
					recv := v.X.(*ast.Ident).Name
					sm := structMethods[recv]
					if sm == nil {
						continue
					}
					for _, m := range sm {
						if m.Name == name {
							m.Dec = fnTyp
							m.Nodes = IfNodesReurned(fnTyp, name)
							break
						}
					}
				}
			}
		}
		return nil
	}, procCtx)

	var buf bytes.Buffer
	buf.WriteString(`// Code generated by script/visitor_gen. DO NOT EDIT.

//go:generate go run github.com/hsiaosiyuan0/mole/script/visitor_gen -d=../parser

package walk

import "github.com/hsiaosiyuan0/mole/ecma/parser"

  `)

	// genereate visitor kindes
	err = genVisitorKinds(&buf, nodeTypStruct, structMethods)
	if err != nil {
		log.Fatal(err)
	}

	// generate visitors
	genVisitors(&buf, nodeTypStruct, structMethods)

	// generate defalt visitors
	err = genDefaultVisitors(&buf, nodeTypStruct)
	if err != nil {
		log.Fatal(err)
	}

	ioutil.WriteFile(distFile, buf.Bytes(), 0644)
	fuzz.Shell("gofmt", "-w", distFile)
}
